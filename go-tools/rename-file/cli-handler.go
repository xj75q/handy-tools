package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
	"sort"
	"time"
)

var (
	handler = NewRename()
	authors = []*cli.Author{
		{
			Name: "coding by qxz",
		},
	}
	cliFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "操作文件路径（必填）",
			Value:    "./",
			Required: true,
		},

		&cli.StringFlag{
			Name:     "output",
			Aliases:  []string{"o"},
			Usage:    "输出文件的路径（选填）",
			Required: false,
		},
	}

	cliCommands = []*cli.Command{
		{
			Name:      "addsign",
			Aliases:   []string{"add"},
			Usage:     "增加文件名标志",
			UsageText: " ",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "addstr",
					Aliases:  []string{"n"},
					Usage:    "默认使用new，如果需要修改为其他名称，请填写",
					Value:    "new",
					Required: false,
				},

				&cli.BoolFlag{
					Name:     "signloc",
					Aliases:  []string{"l"},
					Usage:    "重命名标志位置在左侧，默认true，如需更改使用 -l=false ",
					Value:    true,
					Required: false,
				},
			},
			Action: func(ctx *cli.Context) error {
				handler.addStr = ctx.String("addstr")
				handler.nameLoc = ctx.Bool("signloc")
				handler.outputPath = ctx.String("output")
				inPath := ctx.String("input")
				if inPath == "./" {
					handler.inputPath, _ = os.Getwd()
				} else {
					handler.inputPath = inPath
				}
				handler.InputOperate(ctx)
				time.Sleep(500 * time.Millisecond)
				if err := handler.AddSign(ctx); err != nil {
					return err
				}
				return nil
			},
		},

		{

			Name:      "replace",
			Aliases:   []string{"rep"},
			Usage:     "替换文件的字符串",
			UsageText: " ",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "oldname",
					Aliases:  []string{"v"},
					Usage:    "请填写需替换的字符串",
					Value:    "new",
					Required: false,
				},

				&cli.StringFlag{
					Name:     "newname",
					Aliases:  []string{"n"},
					Usage:    "请填写替换后的字符串",
					Required: false,
				},
			},

			Action: func(ctx *cli.Context) error {
				handler.oldName = ctx.String("oldname")
				handler.newName = ctx.String("newname")
				handler.outputPath = ctx.String("output")

				inPath := ctx.String("input")

				if inPath == "./" {
					handler.inputPath, _ = os.Getwd()
				} else {
					handler.inputPath = inPath
				}
				handler.InputOperate(ctx)
				time.Sleep(500 * time.Millisecond)
				if err := handler.ReplaceFileName(ctx); err != nil {
					return err
				}
				return nil
			},
		},

		{
			Name:      "samename",
			Aliases:   []string{"same"},
			Usage:     "使用相同名字作为文件名",
			UsageText: " ",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "请输入要使用的名字",
					Required: false,
				},
			},
			Action: func(ctx *cli.Context) error {
				sameName := ctx.String("name")
				if sameName == "" {
					log.Println("使用的新名字不能为空，请填写正确的文件名")
					os.Exit(0)
				} else {
					handler.sameFileName = sameName
				}
				handler.outputPath = ctx.String("output")
				inPath := ctx.String("input")
				if inPath == "./" {
					handler.inputPath, _ = os.Getwd()
				} else {
					handler.inputPath = inPath
				}
				handler.InputOperate(ctx)
				time.Sleep(500 * time.Millisecond)

				if err := handler.UseSameName(ctx); err != nil {
					return err
				}
				return nil
			},
		},

		{
			Name:      "usepathname",
			Aliases:   []string{"path"},
			Usage:     "使用所在文件夹名作为新文件的名字",
			UsageText: " ",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:     "frompath",
					Aliases:  []string{"p"},
					Usage:    "是否使用文件夹名字，默认为true，如需关闭设为false",
					Value:    true,
					Required: false,
				},
			},
			Action: func(ctx *cli.Context) error {
				isUsePath := ctx.Bool("p")
				if isUsePath == false {
					return nil
				}
				handler.nameFormPath = isUsePath
				handler.outputPath = ctx.String("output")
				inPath := ctx.String("input")
				if inPath == "./" {
					handler.inputPath, _ = os.Getwd()
				} else {
					handler.inputPath = inPath
				}
				handler.InputOperate(ctx)
				time.Sleep(500 * time.Millisecond)
				if err := handler.UsePathName(ctx); err != nil {
					return err
				}
				return nil
			},
		},

		{
			Name:      "alterSerial",
			Aliases:   []string{"sn"},
			Usage:     "补齐文件前缀的数字",
			UsageText: " ",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:     "serial",
					Aliases:  []string{"s"},
					Usage:    "序号如需补齐2位，填写2；或者填写其他任意位数",
					Value:    2,
					Required: false,
				},
			},
			Action: func(ctx *cli.Context) error {
				serial := ctx.Int("serial")
				if serial < 2 || serial >= 30 {
					log.Println(">> 请输入正确的补齐位数...")
					os.Exit(1)
				}
				handler.filesn = serial
				handler.outputPath = ctx.String("output")
				inPath := ctx.String("input")
				if inPath == "./" {
					handler.inputPath, _ = os.Getwd()
				} else {
					handler.inputPath = inPath
				}
				handler.InputOperate(ctx)
				time.Sleep(500 * time.Millisecond)
				if err := handler.AlterSerial(ctx); err != nil {
					return err
				}
				return nil
			},
		},

		{
			Name:      "substr",
			Aliases:   []string{"sub"},
			Usage:     "删除文件名中的某个字符",
			UsageText: " ",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "subname",
					Aliases:  []string{"n"},
					Usage:    "请填写要删除的字符串",
					Required: false,
				},
				&cli.StringFlag{
					Name:    "subloc",
					Aliases: []string{"l"},
					Usage:   "请填写要删除字符串的位置[默认为全部替换如需更改请使用：left-替换左侧字符，right-替换右侧字符]",

					Value:    "all",
					Required: false,
				},
			},
			Action: func(ctx *cli.Context) error {
				handler.subStr = ctx.String("subname")
				handler.outputPath = ctx.String("output")
				handler.subLoc = ctx.String("subloc")
				inPath := ctx.String("input")
				if inPath == "./" {
					handler.inputPath, _ = os.Getwd()
				} else {
					handler.inputPath = inPath
				}

				handler.InputOperate(ctx)
				time.Sleep(500 * time.Millisecond)
				if err := handler.SubFileName(ctx); err != nil {
					return err
				}
				return nil
			},
		},

		{
			Name:      "rmspace",
			Aliases:   []string{"rms"},
			Usage:     "删除文件名中的空格",
			UsageText: " ",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:     "space",
					Aliases:  []string{"s"},
					Usage:    "是否删除文件名中的空格，默认为true，如需更改设为false",
					Value:    true,
					Required: false,
				},
			},
			Action: func(ctx *cli.Context) error {
				handler.removeSpace = ctx.Bool("space")
				handler.outputPath = ctx.String("output")
				inPath := ctx.String("input")
				if inPath == "./" {
					handler.inputPath, _ = os.Getwd()
				} else {
					handler.inputPath = inPath
				}
				handler.InputOperate(ctx)
				time.Sleep(500 * time.Millisecond)
				if err := handler.RmSpace(ctx); err != nil {
					return err
				}
				return nil
			},
		},
	}
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "【文件批量重命名】"
	app.Usage = "秒级文件批量重命名"
	app.UsageText = "示例：frename -i 文件夹路径 altersn -s 2"
	app.Authors = authors
	app.Flags = cliFlags
	app.Commands = cliCommands
	sort.Sort(cli.FlagsByName(app.Flags))
	defer close(handler.fileInfoChan)
	if err := app.Run(os.Args); err != nil {
		log.Println("\n>> 出错了:", err)
		os.Exit(1)
	}
}
