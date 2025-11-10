package main

import (
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var (
	exePath, _       = os.Executable()
	mailPath         = filepath.Dir(exePath)
	year, month, day = time.Now().Date()
	subject          = fmt.Sprintf("【随手记】 %v-%v-%v", year, fmt.Sprintf("%02d", int(month)), fmt.Sprintf("%02d", day))
)

type (
	CfgInfo struct {
		FromMail string `json:"fromMail"`
		ToMail   string `json:"toMail"`
		Subject  string `json:"subject"`
		Smtp     string `json:"smtp"`
		Pwd      string `json:"pwd"`
	}

	Cfg struct {
		FileName string
		CfgPath  string
		Content  CfgInfo
	}

	MailInfo struct {
		Subject string `json:"subject"`
		Data    string `json:"data"`
	}
)

func mailHandler() *MailInfo {
	return &MailInfo{}
}

func configHandler() *Cfg {
	return &Cfg{
		FileName: "mailcfg.json",
		CfgPath:  mailPath,
	}
}

func (f *Cfg) CreateConfig(fpath string) (error, string) {
	fi := fpath + string(os.PathSeparator) + f.FileName
	file := filepath.Clean(fi)
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		createFile, _ := os.Create(file)
		rb, _ := json.MarshalIndent(f.Content, "", "  ")
		_, err = createFile.Write(rb)
		if err != nil {
			return fmt.Errorf("创建并写入文件失败，请检"), ""
		} else {
			log.Println(">> 邮箱配置文件初始化成功")
			return nil, "success"
		}

	} else {
		var data CfgInfo
		bytes, _ := ioutil.ReadFile(file)
		_ = json.Unmarshal(bytes, &data)
		localCfg := f.Content
		if localCfg.FromMail != "" {
			data.FromMail = localCfg.FromMail
		}
		if localCfg.ToMail != "" {
			data.ToMail = localCfg.ToMail
		}
		if localCfg.Subject != "" {
			data.Subject = localCfg.Subject
		}
		if localCfg.Pwd != "" {
			data.Pwd = localCfg.Pwd
		}
		if localCfg.Smtp != "" {
			data.Smtp = localCfg.Smtp
		}
		result, _ := json.MarshalIndent(data, "", "  ")
		_ = ioutil.WriteFile(file, result, 0644)
		log.Println(">> 邮箱配置文件已更新")
	}

	return nil, ""
}

func (f *Cfg) ReadCfg() (error, *CfgInfo) {
	fi := mailPath + string(os.PathSeparator) + f.FileName
	file := filepath.Clean(fi)
	_, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("邮箱配置文件不存在"), nil
	}
	fileStr := strings.Split(f.FileName, ".")
	viper.SetConfigName(fileStr[0])
	viper.SetConfigType(fileStr[1])
	viper.AddConfigPath(f.CfgPath)

	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("读取配置文件出错：%v\n", err), nil
	}
	cfgInfo := f.Content
	err = viper.Unmarshal(&cfgInfo)
	if err != nil {
		return fmt.Errorf("文件解析出错：%v\n", err), nil
	}
	return nil, &cfgInfo
}

func (c *CfgInfo) IsEmpty() error {
	v := reflect.ValueOf(*c)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		if field.Interface() == reflect.Zero(field.Type()).Interface() {
			if fieldType.Name == "Subject" {
				continue
			} else {
				return fmt.Errorf("配置出错字段出错： '%s' 参数为空，请补充后重试！！！", fieldType.Name)
			}
		}
	}
	return nil
}

func (params *MailInfo) SendMail() error {
	cfg := configHandler()
	err, cfginfo := cfg.ReadCfg()
	if err != nil {
		return err
	}
	if err = cfginfo.IsEmpty(); err != nil {
		return err
	}
	mail := email.NewEmail()
	mail.From = cfginfo.FromMail
	mail.To = []string{cfginfo.ToMail}
	mail.Text = []byte(params.Data)
	if params.Subject != "" {
		mail.Subject = params.Subject
	} else {
		if cfginfo.Subject != "" {
			mail.Subject = cfginfo.Subject
		} else {
			mail.Subject = subject
		}
	}
	addr := cfginfo.Smtp + ":25"
	if err = mail.Send(addr, smtp.PlainAuth("", cfginfo.FromMail, cfginfo.Pwd, cfginfo.Smtp)); err != nil {
		return fmt.Errorf("发送邮件出错:%v", err)
	}
	fmt.Println("send success...")
	return nil
}

var configCommand = &cli.Command{
	Name: "config",
	//Usage:   "Displays global config options and their current values",
	Aliases: []string{"cfg"},

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "from",
			Aliases:  []string{"f"},
			Required: false,
		},

		&cli.StringFlag{
			Name:     "to",
			Aliases:  []string{"t"},
			Required: false,
		},

		&cli.StringFlag{
			Name:     "subj",
			Aliases:  []string{"s"},
			Required: false,
		},

		&cli.StringFlag{
			Name:     "smtp",
			Aliases:  []string{"m"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "pwd",
			Aliases:  []string{"p"},
			Required: false,
		},
	},

	Action: func(c *cli.Context) error {
		cfg := configHandler()
		mail := &CfgInfo{}
		mail.FromMail = c.String("from")
		mail.ToMail = c.String("to")
		mail.Subject = c.String("subj")
		mail.Smtp = c.String("smtp")
		mail.Pwd = c.String("pwd")
		if mail.FromMail == "" || mail.ToMail == "" || mail.Smtp == "" || mail.Pwd == "" {
			return fmt.Errorf("配置参数不能为空，请重新输入")
		}
		if !strings.Contains(mail.FromMail, "@") || !strings.Contains(mail.ToMail, "@") {
			return fmt.Errorf("请输入正确的邮箱参数")
		}
		cfg.Content = *mail
		if err, _ := cfg.CreateConfig(mailPath); err != nil {
			return err
		}
		return nil
	},
}

var mailCommand = &cli.Command{
	Name: "send",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "subj",
			Aliases:  []string{"s"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "data",
			Aliases:  []string{"d"},
			Required: true,
		},
	},

	Action: func(c *cli.Context) error {
		mail := mailHandler()
		mail.Subject = c.String("subj")
		mail.Data = c.String("data")
		if err := mail.SendMail(); err != nil {
			return err
		}
		return nil
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "邮件通知"
	app.HideVersion = true
	app.HideHelpCommand = true
	app.Usage = "(send msg to email...)"
	app.UsageText = fmt.Sprintf("%s\n%s", "./email cfg -f <parames>", "./email send -s <subject> -d <data>")
	app.Commands = []*cli.Command{
		configCommand,
		mailCommand,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Printf(">> %v\n", err)
	}
}
