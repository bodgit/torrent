package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gertd/go-pluralize"
	"github.com/urfave/cli"
	"github.com/zeebo/bencode"
)

// Just enough struct to parse the file list out of a .torrent file
type torrent struct {
	Info struct {
		Files []struct {
			Path []string `bencode:"path"`
		} `bencode:"files"`
		Name string `bencode:"name"`
	} `bencode:"info"`
}

func init() {
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}
}

func clean(c *cli.Context) error {
	if c.NArg() < 2 {
		cli.ShowCommandHelpAndExit(c, c.Command.FullName(), 1)
	}

	f, err := os.Open(c.Args().Get(0))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	var torrent torrent

	dec := bencode.NewDecoder(f)
	if err := dec.Decode(&torrent); err != nil {
		return cli.NewExitError(err, 1)
	}

	if len(torrent.Info.Files) == 0 {
		return cli.NewExitError("single file torrent", 1)
	}

	dir := path.Join(c.Args().Get(1), torrent.Info.Name)

	files := make(map[string]struct{})

	for _, file := range torrent.Info.Files {
		files[path.Join(append([]string{dir}, file.Path...)...)] = struct{}{}
	}

	deleted := 0

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if _, ok := files[path]; !ok {
			if c.Bool("verbose") {
				fmt.Println("Deleting", path)
			}
			if !c.Bool("dry-run") {
				if err := os.Remove(path); err != nil {
					return err
				}
			}
			deleted++
		}
		return nil
	}); err != nil {
		cli.NewExitError(err, 1)
	}

	if deleted > 0 {
		if c.Bool("verbose") {
			client := pluralize.NewClient()
			fmt.Println("Deleted", deleted, "untracked", client.Pluralize("file", deleted, false))
		}
		// Exit with 2 if some files have been deleted
		cli.NewExitError("", 2)
	} else if c.Bool("verbose") {
		fmt.Println("No untracked files")
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
