package main

import (
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/smtp"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var (
	exePath, _       = os.Executable()
	file_path        = filepath.Dir(exePath)
	year, month, day = time.Now().Date()
	subject          = fmt.Sprintf("【随手记】 %v-%v-%v", year, fmt.Sprintf("%02d", int(month)), fmt.Sprintf("%02d", day))
)

type CfgInfo struct {
	FromMail string `json:"fromMail"`
	ToMail   string `json:"toMail"`
	Smtp     string `json:"smtp"`
	Pwd      string `json:"pwd"`
}

type Cfg struct {
	FileName string
	CfgPath  string
	Content  CfgInfo
}

func configHandler() *Cfg {
	return &Cfg{
		FileName: "mail.json",
		CfgPath:  file_path,
	}
}

type Email struct {
	Subject string `json:"subject"`
	Data    string `json:"data"`
}

func mailHandler() *Email {
	return &Email{
		Subject: subject,
	}
}

func StructToMap(obj interface{}) map[string]interface{} {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)

	data := make(map[string]interface{})
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		value := objValue.Field(i).Interface()
		data[field.Name] = value
	}

	return data
}

func (f *Cfg) CreateConfig(fpath string) (error, string) {
	fi := fpath + string(os.PathSeparator) + f.FileName
	file := filepath.Clean(fi)
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		createFile, _ := os.Create(file)
		rb, _ := json.Marshal(f.Content)
		_, err := createFile.Write(rb)
		if err != nil {
			return fmt.Errorf("创建并写入文件失败，请检"), ""
		} else {
			fmt.Println(">> 邮箱配置文件初始化成功")
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

		if localCfg.Pwd != "" {
			data.Pwd = localCfg.Pwd
		}

		if localCfg.Smtp != "" {
			data.Smtp = localCfg.Smtp
		}
		result, _ := json.MarshalIndent(data, "", "")
		_ = ioutil.WriteFile(file, result, 0644)

	}

	return nil, ""
}

func (f *Cfg) ReadCfg() (error, *CfgInfo) {
	fi := file_path + string(os.PathSeparator) + f.FileName
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
	cfg_info := f.Content
	err = viper.Unmarshal(&cfg_info)
	if err != nil {
		return fmt.Errorf("文件解析出错：%v\n", err), nil
	}
	return nil, &cfg_info
}

func (c *CfgInfo) IsEmpty() error {
	v := reflect.ValueOf(*c)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		if field.Interface() == reflect.Zero(field.Type()).Interface() {
			return fmt.Errorf("mail config field '%s' is empty\n", fieldType.Name)
			break
		}
	}
	return nil
}

func (m *Email) SendEmail() error {
	cfg := configHandler()
	err, cfginfo := cfg.ReadCfg()
	if err != nil {
		return err
	}
	if err := cfginfo.IsEmpty(); err != nil {
		return err
	}

	mail := email.NewEmail()
	mail.From = cfginfo.FromMail
	mail.To = []string{cfginfo.ToMail}
	mail.Subject = subject
	mail.Text = []byte(m.Data)
	addr := cfginfo.Smtp + ":25"
	if err := mail.Send(addr, smtp.PlainAuth("", cfginfo.FromMail, cfginfo.Pwd, cfginfo.Smtp)); err != nil {
		return fmt.Errorf("发送邮件出错:%v", err)
	}
	fmt.Println("send success...")
	return nil
}

var configCommand = &cli.Command{
	Name: "config",
	//Usage:   "Displays global config options and their current values",
	Aliases: []string{"c"},

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "from_mail",
			Aliases:  []string{"f"},
			Required: false,
		},

		&cli.StringFlag{
			Name:     "to_mail",
			Aliases:  []string{"t"},
			Required: false,
		},

		&cli.StringFlag{
			Name:     "smtp",
			Aliases:  []string{"s"},
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
		mail.FromMail = c.String("from_mail")
		mail.ToMail = c.String("to_mail")
		mail.Smtp = c.String("smtp")
		mail.Pwd = c.String("pwd")
		cfg.Content = *mail
		cfg.CreateConfig(file_path)
		return nil
	},
}

var mailCommand = &cli.Command{
	Name:    "send",
	Aliases: []string{"s"},

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "data",
			Aliases:  []string{"d"},
			Required: true,
		},
	},

	Action: func(c *cli.Context) error {
		mail := mailHandler()
		mail.Data = c.String("data")
		if err := mail.SendEmail(); err != nil {
			return err
		}
		return nil
	},
}

func main() {
	app := cli.NewApp()
	//app.Name = "发送随笔到邮箱里面"
	app.HideVersion = true
	app.HideHelpCommand = true

	app.Usage = "(send to email...)"
	app.UsageText = `./send c -f <parames> or ./send m -d <parames>`
	app.Commands = []*cli.Command{
		configCommand,
		mailCommand,
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf(">> %v\n", err)
	}
}
