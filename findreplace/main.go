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
		fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
	}

}

func findCommand() *cli.Command {
	return &cli.Command{
		Name:    "find",
		Aliases: []string{"f"},
		Usage:   "find str in stream",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return errors.New("no enough arguments for find command")
			}
			path := c.Args().Get(1)
			search := c.Args().Get(0)
			if len(path) == 0 {
				linesInFile(os.Stdin, search)
			} else {
				err := walkDir(path, search, "")
				if err != nil {
					return err
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
			if c.NArg() == 0 {
				return errors.New("no enough arguments for replace command")
			}
			search := c.Args().Get(0)
			replace := c.Args().Get(1)
			if len(replace) == 0 {
				return errors.New("no enough arguments for replace command")
			}
			path := c.Args().Get(2)
			if len(path) == 0 {
				linesInFile(os.Stdin, search)
			} else {
				err := walkDir(path, search, replace)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func mainAction(arg *cli.Context) error {
	err := arg.App.Command("help").Run(arg)
	if err != nil {
		return errors.Wrap(err, "Help error")
	}
	return nil
}

func findLines(fileName, search string, replace string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open '%s'", fileName)
	}
	var lines []string
	input := bufio.NewScanner(file)
	lineNo := 1
	for input.Scan() {
		text := input.Text()
		if strings.Index(text, search) >= 0 && replace == "" {
			fmt.Printf("%s:%d - %s\n", fileName, lineNo, input.Text())
		} else {
			lines = append(lines, strings.Replace(text, search, replace, -1))
		}
		lineNo++
	}
	defer file.Close()
	return lines, nil

}

func rewriteLines(fileName string, str []string) error {

	file, err := os.Create(fileName)
	if err != nil {
		return errors.Wrap(err, "create file")
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range str {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			return errors.Wrap(err, "write line")
		}
	}
	return w.Flush()
}

func walkDir(path, search, replace string) error {

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		entries, err := ioutil.ReadDir(path)
		if err != nil {
			return errors.Wrap(err, "read dir")
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				lines, err := findLines(filepath.Join(path, entry.Name()), search, replace)
				if err != nil {
					return err
				}
				if len(lines) > 0 {
					err = rewriteLines(filepath.Join(path, entry.Name()), lines)
					if err != nil {
						return err
					}
				}

			}
		}
	case mode.IsRegular():
		lines, err := findLines(path, search, replace)
		if err != nil {
			return err
		}
		if len(lines) > 0 {
			err = rewriteLines(path, lines)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func linesInFile(f *os.File, search string) error {
	input := bufio.NewScanner(f)
	lineNo := 1
	for input.Scan() {
		text := input.Text()
		if strings.Index(text, search) >= 0 {
			fmt.Printf("%d - %s\n", lineNo, input.Text())
		}
		lineNo++
	}
	return nil
}
