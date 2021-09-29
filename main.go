package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dorneanu/tw5-cli/tiddlywiki"
	cli "github.com/urfave/cli/v2"
)

func main() {
	// Define global flags that are always need like the name of the tiddler
	globalFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Value:   "name",
			Usage:   "name of the tiddler",
		},
	}

	var appendText string

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "Get a tiddler",
				Flags:   globalFlags,
				Action: func(c *cli.Context) error {
					tw := tiddlywiki.NewTW(conf.TWHOST)
					_, err := tw.Get("Golang")
					if err != nil {
						fmt.Printf("Errorf: %s", err)
						return err
					}

					// Try to get tiddler
					tiddler, err := tw.Get(c.String("name"))
					if err != nil {
						fmt.Printf("Couldn't get tiddler %s: %s", c.String("name"), err)
						return nil
					}

					fmt.Printf("%s", tiddler.JSON())
					return nil
				},
			},
			{
				Name:    "put",
				Aliases: []string{"p"},
				Flags:   globalFlags,
				Usage:   "put a new tiddler",
				Action: func(c *cli.Context) error {
					fmt.Println("completed task: ", c.Args().First())
					return nil
				},
			},
			{
				Name:    "append",
				Aliases: []string{"a"},
				Flags: append([]cli.Flag{
					&cli.StringFlag{
						Name:        "text",
						Usage:       "Text to append",
						Destination: &appendText,
					},
				}, globalFlags...),
				Usage: "Append content to an existing tiddler",
				Action: func(c *cli.Context) error {
					tw := tiddlywiki.NewTW(conf.TWHOST)
					_, err := tw.Get("Golang")
					if err != nil {
						fmt.Printf("Errorf: %s", err)
						return err
					}

					// Try to get tiddler
					tiddler, err := tw.Get(c.String("name"))
					if err != nil {
						fmt.Printf("Couldn't get tiddler %s: %s", c.String("name"), err)
						return nil
					}

					// Append
					err = tw.Append(tiddler.Title, appendText)
					if err != nil {
						fmt.Printf("Couldn't append to %s", c.String("name"))
						return nil
					}
					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Flags:   globalFlags,
				Usage:   "delete tiddler",
				Action: func(c *cli.Context) error {
					fmt.Println("completed task: ", c.Args().First())
					return nil
				},
			},
		}}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
