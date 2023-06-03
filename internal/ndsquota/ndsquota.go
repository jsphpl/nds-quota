package ndsquota

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/jsphpl/nds-quota/internal/accounts"
	"github.com/jsphpl/nds-quota/internal/renderer"
	"github.com/jsphpl/nds-quota/pkg/ndsctl"
)

type NDSQuota struct {
	accounts *accounts.AccountRepository
	renderer *renderer.Renderer
	ndsctl   *ndsctl.NDSCTL
}

func New(accounts *accounts.AccountRepository, renderer *renderer.Renderer, ndsctl *ndsctl.NDSCTL) *NDSQuota {
	return &NDSQuota{
		accounts: accounts,
		renderer: renderer,
		ndsctl:   ndsctl,
	}
}

// Preauth handles client authentication, rendering an HTML form
// to stdout that is then sent to the user's browser by OpenNDS
func (app *NDSQuota) Preauth(query map[string]string) {
	var params renderer.Params
	params.Query = query

	// NDS is telling us that the user is already authenticated (and probably just
	// returned back to the login page by accident through their browser history)
	if status, ok := query["status"]; ok && status == "authenticated" {
		params.Info = "You are already authenticated."
		app.render(params, false)
		return
	}

	// Someone is trying to authenticate using their token
	if token, ok := query["tk"]; ok {

		// Find the account
		account, err := app.accounts.Find(token)
		if err != nil {
			params.Title = "Invalid Token"
			params.Error = "Please try again"
			app.render(params, true)
			return
		}

		// Ensure the account has still quota remaining
		if account.UsedKiB >= account.QuotaKiB {
			params.Title = "Quota exceeded"
			params.Error = "You have exceeded your quota. Please use a different account."
			app.render(params, true)
			return
		}

		// Authentication successful -> now note down the MAC address and let them in
		client, err := app.ndsctl.ClientInfo(query["clientip"])
		if err != nil {
			params.Title = "Error"
			params.Error = fmt.Sprintf("Unknown error: %s", err)
			app.render(params, true)
			return
		}

		err = app.accounts.AssignMacAddress(account.ID, client.MacAddress)
		if err != nil {
			params.Title = "Error"
			params.Error = fmt.Sprintf("Unknown error: %s", err)
			app.render(params, true)
			return
		}

		err = app.ndsctl.Auth(client.MacAddress)
		if err != nil {
			params.Title = "Error"
			params.Error = fmt.Sprintf("Unknown error: %s", err)
			app.render(params, true)
			return
		}

		// We're done; Send the user a success message and wish them good luck
		params.Title = "Success"
		params.Info = fmt.Sprintf("Have fun browsing the interwebs! Your total quota is %0.2f MiB", float64(account.QuotaKiB)/1024)
		params.Token = client.Token
		app.render(params, false)
		return
	}

	// No other action intended -> render Login page
	params.Title = "Login"
	params.Info = "Welcome! Please enter a valid token in order to access the internet."
	app.render(params, true)
}

func (app *NDSQuota) CheckDeauth() error {
	i := 0
	for account := range app.accounts.All() {
		i += 1
		if len(account.Devices) < 1 {
			// No devices, ergo no quota consumed
			continue
		}

		// Update used quota for account
		used := 0
		for _, device := range account.Devices {
			info, err := app.ndsctl.ClientInfo(device.MacAddress)
			if err != nil {
				return err
			}
			downloaded, err := strconv.Atoi(info.DownloadKiB)
			if err != nil {
				log.Println(err)
				continue
			}
			used += downloaded
		}

		if used > int(account.UsedKiB) {
			account.UsedKiB = uint(used)
			err := app.accounts.Save(account)
			if err != nil {
				log.Println(err)
			}
		}

		// Deauth all of the account's devices if quota exceeded
		if account.UsedKiB >= account.QuotaKiB {
			log.Printf("deauth %s (%v) - used %d/%d KiB", account.ID, account.Devices, account.UsedKiB, account.QuotaKiB)
			for _, device := range account.Devices {
				app.ndsctl.Deauth(device.MacAddress)
			}
		}
	}

	log.Printf("%d accounts scanned", i)

	return nil
}

func (app *NDSQuota) render(params renderer.Params, showForm bool) {
	response, err := app.renderer.Render(params, showForm)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}
}

// ParseQuery process a query-string-like sequence into a map of key-value pairs.
// The input is a string with `key=value` pairs separated by `, ` (comma-space).
func ParseQuery(arg string) (map[string]string, error) {
	urlDec, err := url.QueryUnescape(arg)
	if err != nil {
		return nil, err
	}

	if urlDec[0] == '?' {
		urlDec = urlDec[1:]
	}

	result := make(map[string]string)
	kvToMap(urlDec, &result)

	if fas, ok := result["fas"]; ok {
		fasDec, err := base64.StdEncoding.DecodeString(fas)
		if err != nil {
			return nil, err
		}
		delete(result, "fas")
		kvToMap(string(fasDec), &result)
	}

	return result, nil
}

func kvToMap(input string, m *map[string]string) {
	pairs := strings.Split(input, ", ")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) < 2 {
			continue
		}
		(*m)[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
}
