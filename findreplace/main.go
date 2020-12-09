package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
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
			if c.NArg() == 0 {
				return errors.New("nArgs: no arguments for find command")
			}
			file := c.Args().Get(1)
			str := c.Args().Get(0)
			if len(file) == 0 {
				_, err := findLinesInFile(os.Stdin, str, false)
				if err != nil {
					return errors.Wrap(err, "fndCmdArgFileError:")
				}
			} else {
				err := walkDir(file, str, "")
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
			if c.NArg() == 0 {
				return errors.New("nArgs: no arguments for replace command")
			}
			file := c.Args().Get(2)
			str := c.Args().Get(0)
			replaceStr := c.Args().Get(1)
			if len(file) == 0 {
				err := replaceLinesInFile(os.Stdin, str)
				if err != nil {
					return errors.Wrap(err, "fndCmdArgFileError:")
				}
			} else {
				err := walkDir(file, str, replaceStr)
				if err != nil {
					return errors.Wrap(err, "fndCmdArgWlkError:")
				}
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

func findLines(fileName, str string, strReplace string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "FndLnsOpenError:")
	}
	var lines []string
	input := bufio.NewScanner(f)
	i := 1
	for input.Scan() {
		s := input.Text()
		if strings.Index(s, str) >= 0 && strReplace == "" {
			fmt.Printf("%s:%d - %s\n", fileName, i, input.Text())
		} else {
			lines = append(lines, strings.Replace(s, str, strReplace, -1))
		}
		i++
	}
	f.Close()
	return lines, nil

}

func rewriteLines(fileName string, str []string) error {

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range str {
		fmt.Fprintln(w, line)
	}
	return w.Flush()

}

func walkDir(path, str, replace string) error {

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
				lines, err := findLines(path+entry.Name(), str, replace)
				if err != nil {
					return errors.Wrap(err, "wlkEntriesError:")
				}
				if len(lines) > 0 {
					err = rewriteLines(path+entry.Name(), lines)
					if err != nil {
						return errors.Wrap(err, "wlkRewriteEnriesError:")
					}
				}

			}
		}
	case mode.IsRegular():
		lines, err := findLines(path, str, replace)
		if err != nil {
			return errors.Wrap(err, "wlkRegularError:")
		}
		if len(lines) > 0 {
			err = rewriteLines(path, lines)
			if err != nil {
				return errors.Wrap(err, "wlkRegularRewriteFileError:")
			}
		}

	}
	return nil
}

func findLinesInFile(f *os.File, str string, flagReplace bool) ([]string, error) {
	input := bufio.NewScanner(f)
	i := 1
	for input.Scan() {
		s := input.Text()
		if strings.Index(s, str) >= 0 {
			fmt.Printf("%d - %s\n", i, input.Text())
		}
		i++
		return nil, nil
	}
	return nil, nil
}

func replaceLinesInFile(f *os.File, str string) error {
	input := bufio.NewScanner(f)
	i := 1
	for input.Scan() {
		s := input.Text()
		if strings.Index(s, str) >= 0 {
			fmt.Printf("%d - %s\n", i, input.Text())
		}
		i++
		return nil
	}
	return nil
}
