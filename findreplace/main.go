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
		fmt.Fprintf(os.Stderr, "openError: %v\n", err)
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
				return errors.Wrap(errors.New("NArgs=1"), "No find string in arguments")
			}

			file := os.Args[3]
			if len(file) == 0 {
				findLinesInFile(os.Stdin, os.Args[2])
			} else {
				walkDir(file)
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
	_, err = arg.App.Command("help").Run(arg)
	if err != nil {
		return errors.Wrap(err, "no help yet")
	}
	return nil
}

func findLines(fileName, str string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "can't open file")
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

}

func walkDir(path string) error {
	// todo: добавить возможность подстановки последнего символа '/' в путь файла

	fi, err := os.Stat(path)
	if err != nil {
		return errors.Wrap(err, "no file or directory")
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		entries, err := ioutil.ReadDir(path)
		if err != nil {
			return errors.Wrap(err, "can't read file or directory")
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				_, err = findLines(path+entry.Name(), os.Args[2])
				if err != nil {
					return errors.Wrap(err, "can't find lines")
				}
			}
		}
	case mode.IsRegular():
		_, err = findLines(path, os.Args[2])
		if err != nil {
			return errors.Wrap(err, "can't read file or directory")
		}
	}

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
