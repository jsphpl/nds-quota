package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var ErrAccountNotFound = errors.New("account not found")
var ErrAccountCorrupted = errors.New("account data corrupted")

const (
	accountsDir = "accounts"
)

type Account struct {
	ID       string   `json:"id"`
	QuotaKiB uint     `json:"quota_kib"`
	UsedKiB  uint     `json:"used_kib"`
	Devices  []Device `json:"devices"`
}

type Device struct {
	MacAddress string `json:"id"`
}

type AccountRepository struct {
	location string
}

func NewAccountRepository(location string) *AccountRepository {
	return &AccountRepository{
		location: location,
	}
}

func (repo *AccountRepository) Find(id string) (acc *Account, err error) {
	f, err := os.Open(path.Join(repo.location, accountsDir, id))
	if err != nil {
		return nil, fmt.Errorf("%w for id %s", ErrAccountNotFound, id)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&acc)
	if err != nil {
		return nil, fmt.Errorf("%w for id %s: %w", ErrAccountCorrupted, id, err)
	}

	return
}

func (repo *AccountRepository) AssignMacAddress(accountID, macAddress string) error {
	account, err := repo.Find(accountID)
	if err != nil {
		return err
	}

	for _, device := range account.Devices {
		if device.MacAddress == macAddress {
			// MAC address already assigned to account
			return nil
		}
	}

	account.Devices = append(account.Devices, Device{MacAddress: macAddress})

	return repo.Save(account)
}

func (repo *AccountRepository) Save(account *Account) error {
	f, err := os.OpenFile(path.Join(repo.location, accountsDir, account.ID), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(account)
}

func (repo *AccountRepository) All() <-chan *Account {
	ch := make(chan *Account)

	go func(accounts chan<- *Account) {
		defer close(ch)

		files, err := ioutil.ReadDir(path.Join(repo.location, accountsDir))
		if err != nil {
			return
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			account, err := repo.Find(file.Name())
			if err != nil {
				continue
			}

			ch <- account
		}
	}(ch)

	return ch
}
