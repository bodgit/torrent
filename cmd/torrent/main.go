package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/bodgit/torrent"
	"github.com/gertd/go-pluralize"
	"github.com/urfave/cli"
)

func init() {
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}
}

func readTorrent(file string) (*torrent.Torrent, error) {
	t, err := torrent.New()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := t.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return t, nil
}

func clean(c *cli.Context) error {
	if c.NArg() < 2 {
		cli.ShowCommandHelpAndExit(c, c.Command.FullName(), 1)
	}

	t, err := readTorrent(c.Args().Get(0))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	logger := log.New(ioutil.Discard, "", 0)
	if c.Bool("verbose") {
		logger.SetOutput(os.Stderr)
	}

	deleted, err := t.Clean(c.Args().Get(1), logger, c.Bool("dry-run"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if deleted > 0 {
		client := pluralize.NewClient()
		logger.Println("Deleted", deleted, "untracked", client.Pluralize("file", deleted, false))
		// Exit with 2 if some files have been deleted
		cli.NewExitError("", 2)
	} else {
		logger.Println("No untracked files")
	}

	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "torrent"
	app.Usage = "Utilities for Torrent files"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:  "clean",
			Usage: "Delete any untracked files",
			Description: `Given a Torrent file and its target download directory, delete any files not referenced in the Torrent file.
   Torrent files referencing only one file are not supported.`,
			ArgsUsage: "FILE DIRECTORY",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dry-run, n",
					Usage: "don't actually do anything",
				},
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "increase verbosity",
				},
			},
			Action: clean,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
