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
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	logfile, err := os.OpenFile(config.Logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	fmt.Fprintf(logfile, "preauth call: %+v\n", os.Args)

	renderer, err := renderer.NewRenderer(config.TemplateDirectory)
	if err != nil {
		log.Fatal(err)
	}
	app := ndsquota.New(
		accounts.NewAccountRepository(config.DataDirectory),
		renderer,
		&ndsctl.NDSCTL{},
	)

	if len(os.Args) < 2 {
		log.Fatalf("received no arguments")
	}

	query, err := ndsquota.ParseQuery(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(logfile, "decoded query: %+v\n", query)

	app.Preauth(query)
}
