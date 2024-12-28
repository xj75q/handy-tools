package main

import (
	"fmt"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xpath"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	URL        = "https://www.qweather.com"
	USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 Edg/116.0.1938.69"
)

type weatherData struct {
	city       string
	weather    string
	temp       string
	wind       string
	humidity   string
	airLevel   string
	Aqi        string
	tips       string
	prediction string
}

func weatherHandler() *weatherData {
	return &weatherData{}
}

func (d *weatherData) getWeatherUrl() (error, types.Document) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("User-Agent", USER_AGENT)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Http get err:%v", err), nil

	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Http status code:%v", resp.StatusCode), nil
	}
	defer resp.Body.Close()
	content, err := libxml2.ParseHTMLReader(resp.Body)
	if err != nil {
		return fmt.Errorf("parse error:%v", err), nil

	}
	return nil, content
}

func (d *weatherData) parseData() (error, *weatherData) {
	err, content := d.getWeatherUrl()
	if err != nil {
		return err, nil
	}

	defer content.Free()
	items := xpath.NodeList(content.Find(`//div[@class="current-weather__bg index"]`))
	re := regexp.MustCompile(`\s+`)
	for _, item := range items {
		loc, _ := item.Find(`./div[1]/div[1]/h1/text()`)
		tem, _ := item.Find(`./div[2]/div[@class="current-live__item"]/a/p[1]/text()`)
		weather, _ := item.Find(`./div[2]/div[@class="current-live__item"]/a/p[2]/text()`)
		windLevel, _ := item.Find(`./div[4]/a[1]/p[1]/text()`)
		windDirect, _ := item.Find(`./div[4]/a[1]/p[2]/text()`)
		humidity, _ := item.Find("./div[4]/a[2]/p[1]/text()")
		tip, _ := item.Find(`./div[@class="current-abstract"]/a/text()`)
		pred, _ := item.Find(`./a[@class="live-warning d-flex justify-content-between align-items-center"]/span/text()`)

		d.city = loc.String()
		d.temp = tem.String()
		d.weather = weather.String()
		d.humidity = humidity.String()
		d.prediction = re.ReplaceAllString(pred.String(), " ")
		d.tips = re.ReplaceAllString(re.ReplaceAllString(tip.String(), " "), "")
		result := re.ReplaceAllString(windLevel.String(), " ")
		d.wind = fmt.Sprintf("%v(%v)", windDirect.String(), result)
	}

	aqis := xpath.NodeList(content.Find(`//div[@class="col-6 l-index-left__air"]`))
	for _, aqi := range aqis {
		airlevel, _ := aqi.Find(`./a[@class="c-index-air"]/h3/text()`)
		aqinum, _ := aqi.Find(`./a[@class="c-index-air"]/p/text()`)
		d.airLevel = airlevel.String()
		aqiData := re.ReplaceAllString(re.ReplaceAllString(aqinum.String(), " "), "")
		result := strings.Split(aqiData, "I")[1]
		d.Aqi = result
	}
	return nil, d
}

func main() {
	wh := weatherHandler()
	err, data := wh.parseData()
	if err != nil {
		log.Printf(">> 出错了 %v", err)
		return
	}

	log.Printf("您所在的地区为：%v\n天气情况为    ：%v\n风力          ：%v\n湿度          ：%v\n空气质量      ：%v\nAQI指数为     ：%v\n预计会有      ：%v\n温馨提醒      ：%v\n", data.city, data.weather, data.wind, data.humidity, data.airLevel, data.Aqi, data.prediction, data.tips)
}
