package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/action-workflow-check/internal/commands"
)

var (
	version = "dev"

	cfg struct {
		Scan struct {
			ProjectPath string `arg:"" optional:"" help:"path to project root directory" type:"path" default:"."`
			All         bool   `help:"Enable all rules, otherwise only the action version check is enabled"`
		} `cmd:"" help:"Scan the project for GitHub Actions"`
		Login struct {
		} `cmd:"" help:"Login to GitHub to avoid rate limiting"`
		Debug   bool `help:"Enable debug logging"`
		Version kong.VersionFlag
	}
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cliCtx := kong.Parse(&cfg, kong.Vars{"version": version}, kong.ConfigureHelp(kong.HelpOptions{Compact: true}))
	switch cliCtx.Command() {
	case "scan <project-path>":
		client, err := commands.BuildClient()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to setup github client")
		}

		err = commands.Scan(cliCtx, client, cfg.Scan.ProjectPath, cfg.Debug, cfg.Scan.All)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to scan")
		}
	case "login":
		err := commands.Login()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to login")
		}
	default:
		fmt.Println("Unknown command:", cliCtx.Command())
		os.Exit(1)
	}

}
