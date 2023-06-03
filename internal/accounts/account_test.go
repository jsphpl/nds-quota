package accounts_test

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"testing"

	"github.com/jsphpl/nds-quota/internal/accounts"
)

const location = "./testaccounts"

var users = []accounts.Account{
	{
		ID:       "foo",
		QuotaKiB: 123,
		UsedKiB:  0,
	},
	{
		ID:       "bar",
		QuotaKiB: 100,
		UsedKiB:  99,
	},
	{
		ID:       "broken user",
		QuotaKiB: 100,
		UsedKiB:  99,
	},
}

func TestItFindsAnAccountInFilesystem(t *testing.T) {
	defer deleteTestAccounts()
	createTestAccounts()

	repo := accounts.NewAccountRepository(location)
	u, err := repo.Find("foo")
	if err != nil {
		t.Fatal(err)
	}

	// First user
	if u.ID != "foo" {
		t.Errorf("expected user ID 'foo', got '%s'", u.ID)
	}
	if u.QuotaKiB != 123 {
		t.Errorf("expected quota 123, got %d", u.QuotaKiB)
	}
	if u.UsedKiB != 0 {
		t.Errorf("expected consumed 0, got %d", u.UsedKiB)
	}

	// Second user
	u, err = repo.Find("bar")
	if err != nil {
		t.Fatal(err)
	}
	if u.ID != "bar" {
		t.Errorf("expected user ID 'bar', got '%s'", u.ID)
	}
	if u.QuotaKiB != 100 {
		t.Errorf("expected quota 100, got %d", u.QuotaKiB)
	}
	if u.UsedKiB != 99 {
		t.Errorf("expected consumed 99, got %d", u.UsedKiB)
	}
}

func TestErrsOnMissingAccount(t *testing.T) {
	createTestAccounts()
	defer deleteTestAccounts()

	repo := accounts.NewAccountRepository(location)
	u, err := repo.Find("missing user")
	if !errors.Is(err, accounts.ErrAccountNotFound) {
		t.Errorf("expected error '%s', got '%s'", accounts.ErrAccountNotFound, err)
	}
	if u != nil {
		t.Errorf("expected nil pointer, got %v", u)
	}
}

func TestErrsOnCorrupedJson(t *testing.T) {
	createTestAccounts()
	defer deleteTestAccounts()

	repo := accounts.NewAccountRepository(location)
	u, err := repo.Find("broken user")
	if !errors.Is(err, accounts.ErrAccountCorrupted) {
		t.Errorf("expected error '%s', got '%s'", accounts.ErrAccountNotFound, err)
	}
	if u != nil {
		t.Errorf("expected nil pointer, got %v", u)
	}
}

func createTestAccounts() {
	err := os.Mkdir(location, 6600)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	err = os.Mkdir(path.Join(location, "accounts"), 6600)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	for _, u := range users {
		f, err := os.OpenFile(path.Join(location, "accounts", u.ID), os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		json.NewEncoder(f).Encode(u)

		if u.ID == "broken user" {
			f.WriteAt([]byte("..."), 0)
		}
	}
}

func deleteTestAccounts() {
	os.RemoveAll(location)
}
