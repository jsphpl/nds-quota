package main

import (
	"log"

	"github.com/jsphpl/nds-quota/internal/accounts"
	"github.com/jsphpl/nds-quota/internal/config"
	"github.com/jsphpl/nds-quota/internal/ndsquota"
	"github.com/jsphpl/nds-quota/pkg/ndsctl"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app := ndsquota.New(
		accounts.NewAccountRepository(config.DataDirectory),
		nil,
		&ndsctl.NDSCTL{},
	)

	err = app.CheckDeauth()
	if err != nil {
		log.Fatal(err)
	}
}
