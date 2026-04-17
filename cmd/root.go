package cmd

import (
	"fmt"
	"os"
)

func Execute() error {
	args := os.Args[1:]

	if len(args) == 0 {
		return runSwitch()
	}

	switch args[0] {
	case "new", "create", "add":
		return runNew()
	case "view", "list", "ls":
		return runView()
	case "help", "--help", "-h":
		printHelp()
		return nil
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printHelp()
		return fmt.Errorf("unknown command")
	}
}

func printHelp() {
	fmt.Println(`
  gitswitch — manage git identities per repo

  usage:
    gitswitch           pick a profile and apply it to the current repo
    gitswitch new       create a new profile interactively
    gitswitch view      list all saved profiles
    gitswitch help      show this help

  profiles are stored in ~/.config/git_conf/*.toml
`)
}
