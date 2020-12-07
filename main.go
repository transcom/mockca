package main

import (
	"os"

	"github.com/urfave/cli"
)

func ohnoes(err error) {
	if err != nil {
		panic(err)
	}
}

func Wrapper(cmd func(*cli.Context) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		if err := cmd(c); err != nil {
			panic(err)
		}
		return nil
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "mockca"
	app.Usage = "mockca [command]"
	app.Version = "1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "root", Value: "mockca"},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:   "generate",
			Action: Wrapper(Generate),
			Usage:  "generate",
			Flags: []cli.Flag{
				cli.IntFlag{Name: "bits", Value: 2048},
				cli.IntFlag{Name: "not-before", Value: 0},
				cli.IntFlag{Name: "not-after", Value: 0},
			},
		},
		cli.Command{
			Name:   "sign",
			Action: Wrapper(Sign),
			Usage:  "sign",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "first-name", Value: ""},
				cli.StringFlag{Name: "last-name", Value: ""},
				cli.StringFlag{Name: "middle-name", Value: ""},
				cli.StringFlag{Name: "dod-id", Value: ""},
				cli.StringFlag{Name: "email", Value: ""},
				cli.StringFlag{Name: "org", Value: ""},
			},
		},
	}

	app.Run(os.Args)
}
