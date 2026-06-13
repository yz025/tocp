package main

import (
	"context"
	"log"
	"os"

	cp "github.com/otiai10/copy"
	"github.com/urfave/cli/v3"
)

func RunTOCP() {
	log.SetFlags(0)

	cmd := &cli.Command{
		Name:  "tocp",
		Usage: "TOML Copy Paste",
		Commands: []*cli.Command{
			{
				Name:  "push",
				Usage: "Copy src to dst",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return run(cmd, true)
				},
			},
			{
				Name:  "pull",
				Usage: "Copy dst to src",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return run(cmd, false)
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "path to a tocp.toml",
				Required: false,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logError(err)
		os.Exit(1)
	}
}
func run(cmd *cli.Command, isPush bool) error {
	pathFlag := cmd.String("path")

	config, err := LoadConfig(pathFlag)
	if err != nil {
		return err
	}

	succeeded := 0
	for _, pair := range config.Pairs {
		if update(pair, isPush) {
			succeeded += 1
		}

		log.Println()
	}
	log.Printf("%d/%d succeeded.", succeeded, len(config.Pairs))

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
		logFailuref(err, "Failed to find %s %q.", fromStr, from)
		return false
	}

	if *pair.RemoveBeforeCopy {
		if err := os.RemoveAll(to); err != nil {
			logFailuref(err, "Failed to remove %s %q.", toStr, to)
			return false
		} else {
			logSuccessf("Removed %s %q", toStr, to)
		}
	}

	if err := cp.Copy(from, to); err != nil {
		logFailuref(err, "Failed to copy %s %q -> %s %q.", fromStr, from, toStr, to)
		return false
	} else {
		logSuccessf("Copied %s %q -> %s %q", fromStr, from, toStr, to)
	}

	return true
}
func logError(err error) {
	log.Printf("\x1b[31mError\x1b[0m: %v", err)
}
func logSuccessf(format string, v ...any) {
	str := "\x1b[32mO\x1b[0m " + format
	log.Printf(str, v...)
}
func logFailuref(err error, format string, v ...any) {
	str := "\x1b[31mX\x1b[0m " + format
	log.Printf(str, v...)
	logError(err)
}
