package initcmd

import (
	"flag"
	"fmt"
)

type InitFlags struct {
	Name  string
	Email string
}

var Flags InitFlags

func RunInit(args []string) error {
	initFlags := flag.NewFlagSet("init", flag.ContinueOnError)
	initFlags.StringVar(&Flags.Email, "e", "", "Alias of --email")
	initFlags.StringVar(&Flags.Email, "email", "", "Your email for userAgent header")

	initFlags.StringVar(&Flags.Name, "n", "", "Alias of --name")
	initFlags.StringVar(&Flags.Name, "name", "", "Your name for userAgent header")

	initFlagsErr := initFlags.Parse(args)
	if initFlagsErr != nil {
		return initFlagsErr
	}

	initFlags.Parse(args)

	rest := initFlags.Args()
	fmt.Println(rest)
	fmt.Println(Flags.Email)
	fmt.Println(Flags.Name)

	return nil
}
