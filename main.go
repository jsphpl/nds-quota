package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jsphpl/nds-quota/internal/accounts"
	"github.com/jsphpl/nds-quota/internal/config"
	"github.com/jsphpl/nds-quota/internal/ndsquota"
	"github.com/jsphpl/nds-quota/internal/renderer"
	"github.com/jsphpl/nds-quota/pkg/ndsctl"
)

func main() {
	// load config and open logfile
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	logfile, err := os.OpenFile(config.Logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()

	// Build the application with all its dependencies
	renderer, err := renderer.NewRenderer(config.TemplateDirectory)
	if err != nil {
		log.Fatal(err)
	}
	app := ndsquota.New(
		accounts.NewAccountRepository(config.DataDirectory),
		renderer,
		&ndsctl.NDSCTL{},
	)

	// Run it!
	if len(os.Args) < 2 {
		fmt.Fprintf(logfile, "missing argument")
	} else if os.Args[1] == "check-deauth" {
		fmt.Fprintf(logfile, "checkin quota")
		err = app.CheckDeauth()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		query, err := ndsquota.ParseQuery(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(logfile, "preauth call: %+v\n", query)
		app.Preauth(query)
	}
}
