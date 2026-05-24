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
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(io.Discard)
	}

	for _, pair := range config.Pairs {
		src := os.Expand(pair.Src, os.Getenv)
		dst := os.Expand(pair.Dst, os.Getenv)

		if isPush {
			update(src, dst)
		} else {
			update(dst, src)
		}
	}
	return nil
}
func update(src string, dst string) {
	_, err := os.Stat(src)
	if err != nil {
		log.Printf("WARN: %s NOT FOUND\n", src)
		return
	}

	if err := os.RemoveAll(dst); err != nil {
		log.Printf("ERR: FAILED TO REMOVE %s / %s\n", dst, err)
	} else {
		log.Printf("LOG: REMOVED %s\n", dst)
	}

	if err := cp.Copy(src, dst); err != nil {
		log.Printf("ERR: FAILED TO COPY %s TO %s / %s\n", src, dst, err)
		return
	}
	log.Printf("LOG: COPIED %s TO %s\n", src, dst)
}
