package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	cp "github.com/otiai10/copy"
	"github.com/urfave/cli/v3"
)

func RunTOCP() {
	cmd := &cli.Command{
		Name:  "tocp",
		Usage: "TOML Copy Paste",
		Commands: []*cli.Command{
			{
				Name:  "push",
				Usage: "Copy `src` to `dst`.",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return run(cmd, true)
				},
			},
			{
				Name:  "pull",
				Usage: "Copy `dst` to `src`",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return run(cmd, false)
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Path to a tocp.toml",
				Required: false,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
func run(cmd *cli.Command, isPush bool) error {
	pathFlag := cmd.String("path")

	config, err := LoadConfig(pathFlag)
	if err != nil {
		return err
	}

	if config.Log {
		log.SetFlags(0)
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(io.Discard)
	}

	succeeded := 0
	failed := 0
	for _, pair := range config.Pairs {
		if update(pair, isPush) {
			succeeded += 1
		} else {
			failed += 1
		}

		log.Println()
	}

	fmt.Printf("\033[32m%d succeeded\033[0m, \033[31m%d failed\033[0m\n", succeeded, failed)

	return nil
}
func update(pair Pair, isPush bool) bool {
	var from, fromStr string
	var to, toStr string

	if isPush {
		from = os.Expand(pair.Src, os.Getenv)
		fromStr = "src"

		to = os.Expand(pair.Dst, os.Getenv)
		toStr = "dst"
	} else {
		from = os.Expand(pair.Dst, os.Getenv)
		fromStr = "dst"

		to = os.Expand(pair.Src, os.Getenv)
		toStr = "src"
	}

	_, err := os.Stat(from)
	if err != nil {
		log.Printf("\033[31mX\033[0m Failed to find %s\n", from)
		return false
	}

	if *pair.RemoveBeforeCopy {
		if err := os.RemoveAll(to); err != nil {
			log.Printf("\033[31mX\033[0m Failed to remove %s %s / %s\n", toStr, to, err)
			return false
		} else {
			log.Printf("\033[32mO\033[0m Removed %s %s\n", toStr, to)
		}
	}

	if err := cp.Copy(from, to); err != nil {
		log.Printf("\033[31mX\033[0m Failed to copy %s %s -> %s %s / %s\n", fromStr, from, toStr, to, err)
		return false
	} else {
		log.Printf("\033[32mO\033[0m Copied %s %s -> %s %s\n", fromStr, from, toStr, to)
	}
	return true
}
