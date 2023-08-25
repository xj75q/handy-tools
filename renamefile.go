package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type FieldName struct {
	InputPath    string
	OutputPath   string
	NameFormPath bool
	NewName      string
	Type         string
}

func NewHandler() *FieldName {
	return &FieldName{}
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func (f *FieldName) RenameFiles() error {
	suffixFlag := f.JudgeFileType()
	var newfile string
	err := filepath.Walk(f.InputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if IsFile(path) && strings.HasSuffix(path, suffixFlag) {
			oldfile := path
			sList := strings.Split(oldfile, "/")
			originFilePath := strings.Join(sList[:len(sList)-1], "/")
			var build strings.Builder
			build.WriteString(f.NewName)
			build.WriteString("[a-zA-Z0-9_\u4e00-\u9fa5]+")
			nameStr := build.String()
			reg, _ := regexp.Compile(nameStr)
			newName := reg.FindString(originFilePath)
			if f.OutputPath == "" {
				newfile = originFilePath + "/" + newName + suffixFlag
			} else if strings.HasSuffix(f.OutputPath, "/") {
				newfile = f.OutputPath + newName + suffixFlag
			} else {
				newfile = f.OutputPath + "/" + newName + suffixFlag
			}

			fmt.Println(oldfile, newfile)
			err := os.Rename(oldfile, newfile)
			if err != nil {
				return err
			}
			fmt.Println("rename files success...")

		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (f *FieldName) JudgeFileType() string {
	switch f.Type {
	case "go":
		return ".go"
	case "py":
		return ".py"
	case "log":
		return ".log"
	case "text":
		return ".text"
	case "docx":
		return ".docx"
	case "pptx":
		return ".pptx"
	default:
		return ".go"
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "rename files program"
	app.Usage = "input the file path and auto switch to new file"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "please input file path",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "output",
			Aliases:  []string{"o"},
			Usage:    "please input the output path ",
			Required: false,
		},
		&cli.StringFlag{
			Name:        "name",
			Aliases:     []string{"n"},
			Usage:       "please input the new name,if you take a value from path,then input the name key worlds",
			DefaultText: "new",
			Required:    true,
		},
		&cli.BoolFlag{
			Name:     "frompath",
			Aliases:  []string{"f"},
			Usage:    "Whether to take a value from the input file path,default 'false' ",
			Value:    false,
			Required: false,
		},
		&cli.StringFlag{
			Name:        "type",
			Aliases:     []string{"t"},
			Usage:       "please input the file type,eg: go, python ,txt",
			DefaultText: "go",
			Required:    false,
		},
	}

	app.Action = func(c *cli.Context) error {
		handler := NewHandler()
		handler.InputPath = c.String("input")
		handler.OutputPath = c.String("output")
		handler.NameFormPath = c.Bool("frompath")
		handler.NewName = c.String("name")
		handler.Type = c.String("type")
		err := handler.RenameFiles()
		if err != nil {
			return err
		}
		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(">> there is an err:", err)
	}
}
