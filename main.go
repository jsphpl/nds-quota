package main

import (
	"log"
	"os"
	"sync"

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

	var logpath string
	if config.Log.Enabled {
		logpath = config.Log.Path
	} else {
		logpath = os.DevNull
	}
	logfile, err := os.OpenFile(logpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	// Build the application with all its dependencies
	renderer, err := renderer.NewRenderer(config.TemplateDirectory)
	if err != nil {
		log.Fatal(err)
	}
	app := ndsquota.New(
		accounts.NewAccountRepository(config.DataDirectory),
		renderer,
		ndsctl.NewNDSCTL(config.NDSCTLBin, &sync.Mutex{}),
	)

	// Run it!
	if len(os.Args) < 2 {
		log.Printf("missing argument")
	} else if os.Args[1] == "check-deauth" {
		log.Printf("checking quota")
		err = app.CheckDeauth()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		query, err := ndsquota.ParseQuery(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("preauth call: %+v\n", query)
		app.Preauth(query)
	}
}
