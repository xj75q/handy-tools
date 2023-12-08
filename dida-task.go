package main

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var (
	BaseUrl       = "https://www.dida365.com"
	TaskApiUrl    = BaseUrl + "/api/v2/task"
	ProjectApiUrl = BaseUrl + "/api/v2/projects"
	LoginUrl      = BaseUrl + "/api/v2/user/signon?wc=true&remember=true"
	exe_path, _   = os.Executable()
	fpath         = filepath.Dir(exe_path)
	fname         = "dida-cfg.json"
)

type userInfo struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type cfgInfo struct {
	Cookie      string `json:"cookie"`
	ProjectName string `json:"projectname"`
	ProjectId   string `json:"projectId"`
}

type cfg struct {
	FileName string
	CfgPath  string
	Content  cfgInfo
}

type htmlParams struct {
	Data   string
	Method string
}

func cfgHandler() *cfg {
	return &cfg{
		FileName: fname,
		CfgPath:  fpath,
	}
}

func userHandler() *userInfo {
	return &userInfo{}
}

func htmlHandler() *htmlParams {
	return &htmlParams{
		Method: "POST",
	}
}

func (c *cfg) initConfig() (error, string) {
	file := filepath.Clean(fpath + string(os.PathSeparator) + fname)
	fmt.Println(file)
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			createFile, _ := os.Create(file)
			rb, _ := json.Marshal(c.Content)
			_, err := createFile.Write(rb)
			if err != nil {
				return fmt.Errorf("创建并写入文件失败，请检查..."), ""
			} else {
				fmt.Println(">> 配置文件初始化成功")
				return nil, "success"
			}
		} else {
			fmt.Printf("无法判断文件 %s 是否存在：%v\n", file, err)
		}
	} else {
		result, _ := json.MarshalIndent(c.Content, "", "")
		_ = ioutil.WriteFile(file, result, 0644)
		return nil, "update config file success..."
	}
	return nil, ""
}

func (c *cfg) readCfg() (error, *cfgInfo) {
	file := filepath.Clean(fpath + string(os.PathSeparator) + c.FileName)
	_, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("配置文件不存在"), nil
	}

	fileStr := strings.Split(c.FileName, ".")
	viper.SetConfigName(fileStr[0])
	viper.SetConfigType(fileStr[1])
	viper.AddConfigPath(c.CfgPath)

	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("读取配置文件出错：%v\n", err), nil
	}
	cfginfo := c.Content
	err = viper.Unmarshal(&cfginfo)
	if err != nil {
		return fmt.Errorf("文件解析出错：%v\n", err), nil
	}
	return nil, &cfginfo
}

func (c *cfgInfo) isEmpty() error {
	v := reflect.ValueOf(*c)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		if field.Interface() == reflect.Zero(field.Type()).Interface() {
			return fmt.Errorf("dida config field '%s' is empty\n", fieldType.Name)
			break
		}
	}
	return nil
}

func cfgData() (error, *cfg) {
	cfg := cfgHandler()
	err, cfginfo := cfg.readCfg()
	if err != nil {
		return err, nil
	}
	if err := cfginfo.isEmpty(); err != nil {
		return err, nil
	}
	return nil, cfg
}

func (h *htmlParams) generateReqHeader(cookie string) map[string]interface{} {
	var info interface{}
	headerStr := `
				  {"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:36.0) Gecko/20100101 Firefox/36.0",
                  "Accept-Language": "zh-CN,en-US;q=0.7,en;q=0.3", 
	              "DNT": "1",
	              "Accept": "application/json, text/javascript, */*; q=0.01",
	              "Content-Type": "application/json; charset=UTF-8",
	              "X-Requested-With": "XMLHttpRequest",
	              "Accept-Encoding": "deflate"}
	`
	if err := json.Unmarshal([]byte(headerStr), &info); err != nil {
		fmt.Errorf("")
	}
	headers := info.(map[string]interface{})
	if cookie != "" {
		headers["Cookie"] = cookie
	}
	headers["Referer"] = BaseUrl
	return headers
}

func (u *userInfo) login() error {
	web := htmlHandler()
	client := &http.Client{}
	sendData, _ := json.Marshal(&u)
	req, err := http.NewRequest(web.Method, LoginUrl, strings.NewReader(string(sendData)))
	if err != nil {
		return fmt.Errorf("")
	}
	headers := web.generateReqHeader("")
	for key, header := range headers {
		req.Header.Set(key, header.(string))
	}
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
	cookie := "t=" + respData["token"].(string)
	cfg := cfgHandler()
	cfg.Content.Cookie = cookie
	cfg.initConfig()

	return nil

}

func (c *cfg) checkLogin() (error, *cfgInfo) {
	err, data := c.readCfg()
	if err != nil {
		return err, nil
	}
	if data.Cookie == "" {
		return fmt.Errorf("请先登录..."), nil
	}
	return nil, data
}

func (c *cfg) checkProject() (error, *cfgInfo) {
	err, data := c.checkLogin()
	if err != nil {
		return err, nil
	}
	if data.ProjectId == "" {
		return fmt.Errorf("请先设置清单名"), nil
	}
	return nil, data
}

func (c *cfg) setProject() error {
	web := htmlHandler()
	err, info := c.checkLogin()
	if err != nil {
		return err
	}
	c.Content.Cookie = info.Cookie
	client := &http.Client{}
	req, err := http.NewRequest("GET", ProjectApiUrl, strings.NewReader(""))
	if err != nil {
		return err
	}
	headers := web.generateReqHeader(info.Cookie)
	for key, header := range headers {
		req.Header.Set(key, header.(string))
	}
	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resp status:%v", resp.Status)
	}
	defer resp.Body.Close()

	var respData = []map[string]interface{}{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	jErr := json.Unmarshal(body, &respData)
	if jErr != nil {
		return jErr
	}

	for _, project := range respData {
		inputName := c.Content.ProjectName
		if project["name"].(string) == inputName {
			c.Content.ProjectId = project["id"].(string)
			//todo
			err, result := c.initConfig()
			if err != nil {
				return err
			}
			fmt.Println(result)
			break
		}
	}
	return nil
}

func recordText(title, content, projectId string) map[string]interface{} {
	record := make(map[string]interface{})
	record["modifiedTime"] = time.Now().Format("2006-01-02 15:04:05")
	record["title"] = title
	record["priority"] = 0
	record["status"] = 0
	record["deleted"] = 0
	record["content"] = content
	record["sortOrder"] = 0
	record["projectId"] = projectId
	record["startDate"] = ""
	record["dueDate"] = ""
	record["items"] = []string{}
	record["assignee"] = ""
	record["progress"] = 0
	record["tags"] = []string{}
	record["isAllDay"] = true
	record["reminder"] = ""
	record["local"] = true
	record["isDirty"] = false
	return record
}

func (c *cfg) sendTask(title, content string) error {
	err, LocalCfg := c.checkProject()
	if err != nil {
		fmt.Println("")
	}

	data := recordText(title, content, LocalCfg.ProjectId)
	web := htmlHandler()
	client := &http.Client{}
	sendData, _ := json.Marshal(&data)
	req, err := http.NewRequest(web.Method, TaskApiUrl, strings.NewReader(string(sendData)))
	if err != nil {
		return fmt.Errorf("")
	}
	headers := web.generateReqHeader(LocalCfg.Cookie)
	for key, header := range headers {
		req.Header.Set(key, header.(string))
	}
	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resp status:%v", resp.Status)
	}
	defer resp.Body.Close()

	fmt.Println(">> 记录清单完成 ...")
	return nil
}

var (
	didaCmd = &cobra.Command{
		Short: "命令行创建任务清单到（滴答清单）",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	registerCmd = &cobra.Command{
		Use:     "login",
		Short:   "登录滴答清单",
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			user := userHandler()
			user.UserName, _ = cmd.Flags().GetString("username")
			user.Password, _ = cmd.Flags().GetString("password")
			if err := user.login(); err != nil {
				fmt.Println(err)
			}
		},
	}

	projectCmd = &cobra.Command{
		Use:   "project",
		Short: "设置清单目录",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := cfgHandler()
			input, _ := cmd.Flags().GetString("name")
			cfg.Content.ProjectName = input
			if err := cfg.setProject(); err != nil {
				fmt.Println(err)
			}
		},
	}

	taskCmd = &cobra.Command{
		Use:     "task",
		Short:   "创建任务",
		Aliases: []string{},
		Run: func(cmd *cobra.Command, args []string) {
			cfg := cfgHandler()
			title, _ := cmd.Flags().GetString("title")
			content, _ := cmd.Flags().GetString("content")
			if err := cfg.sendTask(title, content); err != nil {
				fmt.Println(err)
			}

		},
	}
)

func init() {
	registerCmd.Flags().StringP("username", "u", "", "用户名")
	registerCmd.Flags().StringP("password", "p", "", "密码")
	projectCmd.Flags().StringP("name", "n", "", "清单名")
	taskCmd.Flags().StringP("title", "i", "", "任务标题")
	taskCmd.Flags().StringP("content", "t", "", "任务内容")
	registerCmd.MarkFlagsRequiredTogether("username", "password")

	didaCmd.AddCommand(registerCmd)
	didaCmd.AddCommand(projectCmd)
	didaCmd.AddCommand(taskCmd)
}

func main() {
	if err := didaCmd.Execute(); err != nil {
		fmt.Printf("err of %s", err)
		os.Exit(1)
	}
}
