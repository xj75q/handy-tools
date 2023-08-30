package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	Conf        ConfigFile
	filePath, _ = os.Getwd()
)

type HtmlParams struct {
	Data   string
	Method string
}

type ConfigFile struct {
	FileName string
	CfgPath  string
	FlomoApi string
}

func NewHtmlHandler() *HtmlParams {
	return &HtmlParams{
		Method: "POST",
	}
}

func NewFileHandler() *ConfigFile {
	return &ConfigFile{
		FileName: "flomoCfg.json",
		CfgPath:  filePath,
	}
}

// 查看文件是否存在（使用os.Stat()函数判断文件或文件夹是否存在）
func (f *ConfigFile) isFileExist(fpath string) (error, string) {
	fi := fpath + "/" + f.FileName
	file := filepath.Clean(fi)

	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			createFile, _ := os.Create(file)
			jsonStr := `{"flomoApi":""}`
			_, err := createFile.WriteString(jsonStr)
			if err != nil {
				return fmt.Errorf("创建并写入文件失败，请检查目录权限"), ""
			} else {
				fmt.Println(">> flomo配置文件初始化成功，请填入你的API")
				return nil, "success"
			}
		} else {
			fmt.Printf("无法判断文件 %s 是否存在：%v\n", file, err)
		}
	} else {
		err, url := f.ReadConfig()
		if err != nil {
			return err, ""
		} else {
			return nil, url
		}
	}
	return nil, ""

}

// viper读取配置文件中的内容
func (f *ConfigFile) ReadConfig() (error, string) {
	fileStr := strings.Split(f.FileName, ".")
	viper.SetConfigName(fileStr[0])
	viper.SetConfigType(fileStr[1])
	viper.AddConfigPath(f.CfgPath)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("读取配置文件出错：%v\n", err), ""
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		return fmt.Errorf("文件解析出错：%v\n", err), ""
	}

	isFlomoKey := strings.Contains(Conf.FlomoApi, "https://flomoapp.com/iwh")
	if !isFlomoKey {
		return fmt.Errorf("flomo Key不正确，请在配置文件中填入正确的apiKey"), ""
	}
	return nil, Conf.FlomoApi
}

func (h *HtmlParams) SendPost(url, data string) error {
	client := &http.Client{}
	postData := make(map[string]string)
	postData["content"] = data
	sendData, _ := json.Marshal(postData)
	req, err := http.NewRequest(h.Method, url, strings.NewReader(string(sendData)))
	if err != nil {
		return fmt.Errorf("发送到flomo笔记出错:%v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resp status:%v", resp.Status)
	}
	defer resp.Body.Close()
	var respData = make(map[string]interface{})
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jErr := json.Unmarshal(body, &respData)
	if jErr != nil {
		return jErr
	}
	codeFloat := respData["code"]
	codeStr := fmt.Sprintf("%f", codeFloat)
	code := strings.Split(codeStr, ".")[0]
	if code == "-1" {
		return fmt.Errorf("请使用flomo会员，才能发送数据...")
	} else {
		fmt.Printf("已发送到flomo的数据为：%v", data)
		return nil
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "发送随笔到flomo笔记"
	app.Usage = "(send to flomo...)"
	app.UsageText = `./send -i 'data-1'`
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Required: true,
		},
	}

	app.Action = func(c *cli.Context) error {
		data := c.String("input")
		fileHandler := NewFileHandler()
		htmlHandler := NewHtmlHandler()
		err, url := fileHandler.isFileExist(filePath)
		if err != nil {
			return err
		}
		if url == "success" {
			return nil
		}
		err = htmlHandler.SendPost(url, data)

		return err
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf(">> %v\n", err)
	}
}
