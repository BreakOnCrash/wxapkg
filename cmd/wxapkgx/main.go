package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ac0d3r/wxapkg"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	err_  = color.New(color.FgRed)
	warn_ = color.New(color.FgYellow)
	info_ = color.New(color.FgCyan)
)

func New() *cli.App {
	app := cli.NewApp()
	app.Name = "wxapkg"
	app.Usage = "wxapkg analysis tool for macos"
	app.Version = "v0.0.2"
	app.Commands = []*cli.Command{
		{
			Name:  "unpack",
			Usage: "unpack .wxapkg file",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "in", Required: true, Usage: ".wxapkg file path"},
				&cli.StringFlag{Name: "out", Required: false, Value: "./unpack_out", Usage: "unpacked output path"},
				&cli.BoolFlag{Name: "format", Value: false, Usage: "format content (e.g. js html json)"},
				&cli.BoolFlag{Name: "v", Value: false, Usage: "more info"},
			},
			Action: unpack,
		},
		{
			Name:   "list",
			Usage:  "list macOS Wechat .wxapkg file",
			Flags:  []cli.Flag{},
			Action: list,
		},
		{
			Name:  "debug",
			Usage: "start debug webview (TODO)",
		},
	}

	return app
}

func unpack(c *cli.Context) error {
	in := c.String("in")
	out := c.String("out")
	v := c.Bool("v")
	format := c.Bool("format")

	info_.Printf("[+] unpacking %s -> %s \n", in, out)
	var infof func(format string, a ...interface{}) = nil
	if v {
		infof = func(format string, a ...interface{}) {
			info_.Printf(format, a...)
		}
	}

	err := wxapkg.Unpack(in, out, format, infof)
	if err != nil {
		if errors.Is(err, wxapkg.ErrInvalidWXAPkg) {
			err_.Printf("[-] '%s' %s\n", in, err.Error())
		}
		err_.Printf("[-] unpack '%s' error: %s \n", in, err.Error())
	}

	info_.Println("[+] unpacking completed")
	return nil
}

func list(c *cli.Context) error {
	warn_.Println("[*] only support WeChat version 3.8.*")
	root := wxapkg.GetWXAppletPath()
	info_.Printf("list %s\n", root)

	stat, err := os.Stat(root)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		err_.Println("[-] not support WeChat version")
		return nil
	}

	var tab int
	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path == root {
				return nil
			}

			tab++
			info_.Printf("%s%s\n", getTabTree(tab), path)
		} else {
			info_.Printf("%s%s\n", getTabTree(tab), path)
		}

		return nil
	})
}

func getTabTree(tab int) string {
	s := ""
	for i := 0; i < tab; i++ {
		s += "\t"
	}
	return s
}

func main() {
	if err := New().Run(os.Args); err != nil {
		err_.Println()
		err_.Printf("[-]  %s\n", err.Error())
	}
}
