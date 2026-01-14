package params

import (
	"log"
	initcmd "market-data/src/init"

	"os"
)

func ParamsResover() {
	args := os.Args

	if len(args) < 2 {
		log.Fatalf("Expected subcommand.")
	}

	switch args[1] {
	case "init":
		runInitErr := initcmd.RunInit(args[2:])
		if runInitErr != nil {
			log.Fatalf("ParamsResover: init error %v\n%v\n", runInitErr, os.Stderr)
		}
	case "update":
		runInitErr := RunUpdate(args[2:])
		if runInitErr != nil {
			log.Fatalf("ParamsResover: init error %v\n%v\n", runInitErr, os.Stderr)
		}
	default:
		log.Fatalf("Unknown subcommand: %v\n", args[1])
	}
}
