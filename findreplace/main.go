package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "findreplace"
	app.Usage = "Find and replace string in file(s) or stdin"
	app.Description = "Description"
	app.Authors = []*cli.Author{
		{Name: "Andrew Gly", Email: "glybin.av@gmail.com"},
	}

	app.Commands = []*cli.Command{
		findCommand(),
		replaceCommand(),
	}

	app.Action = mainAction

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runError: %v\n", err)
	}

}

func findCommand() *cli.Command {
	return &cli.Command{
		Name:    "find",
		Aliases: []string{"f"},
		Usage:   "find str in stream",
		Action: func(c *cli.Context) error {
			n := c.NArg()
			if n == 0 {
				return errors.Wrap(errors.New("NArgs"), "no arguments for find command")
			}
			if n == 1 {
				return errors.Wrap(errors.New("NArgs=1"), "no find string in arguments")
			}

			file := os.Args[3]
			if len(file) == 0 {
				err := findLinesInFile(os.Stdin, os.Args[2])
				if err != nil {
					return errors.Wrap(err, "fndCmdArgFileError:")
				}
			} else {
				err := walkDir(file)
				if err != nil {
					return errors.Wrap(err, "fndCmdArgWlkError:")
				}
			}
			return nil
		},
	}
}

func replaceCommand() *cli.Command {
	return &cli.Command{
		Name:    "replace",
		Aliases: []string{"r"},
		Usage:   "replace str in stream",
		Action: func(c *cli.Context) error {
			n := c.NArg()
			if n == 0 {
				return errors.Wrap(errors.New("NArgs"), "no arguments for find command")
			}

			return nil
		},
	}
}

func mainAction(arg *cli.Context) error {
	err := arg.App.Command("help").Run(arg)
	if err != nil {
		return errors.Wrap(err, "HlpError:")
	}
	return nil
}

func findLines(fileName, str string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "FndLnsOpenError:")
	}
	input := bufio.NewScanner(f)
	i := 1
	for input.Scan() {
		s := input.Text()
		if strings.Index(s, str) >= 0 {
			fmt.Printf("%s:%d - %s\n", fileName, i, input.Text())
		}
		i++
	}
	f.Close()
	return nil

}

func walkDir(path string) error {
	// todo: добавить возможность подстановки последнего символа '/' в путь файла

	fi, err := os.Stat(path)
	if err != nil {
		return errors.Wrap(err, "wlkPathError:")
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		entries, err := ioutil.ReadDir(path)
		if err != nil {
			return errors.Wrap(err, "wlkReadDirError:")
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				err = findLines(path+entry.Name(), os.Args[2])
				if err != nil {
					return errors.Wrap(err, "wlkEntriesError:")
				}
			}
		}
	case mode.IsRegular():
		err = findLines(path, os.Args[2])
		if err != nil {
			return errors.Wrap(err, "wlkRegularError:")
		}
	}
	return nil
}

func findLinesInFile(f *os.File, str string) error {
	input := bufio.NewScanner(f)
	i := 1
	for input.Scan() {
		s := input.Text()
		if strings.Index(s, str) >= 0 {
			fmt.Printf("%s:%d - %s\n", filepath.Base(os.Args[3]), i, input.Text())
		}
		i++
	}
	return nil
}
