package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	stream       = make(chan interface{})
	finalStream  = make(chan string)
	minusHour, _ = time.ParseDuration("-1h")
	didaNow      = time.Now().Add(8 * minusHour) //8小时的时差
	midday       = 12 - 8                        //需要在12点的基础上减去8小时
	defaultNight = "13:00:00"                    //晚上21点，减去8小时后就是现在设置的13：00
	timeFlag     = ".000+0000"
)

type TimeModify struct {
}

func TimeHandler() *TimeModify {
	return &TimeModify{}
}

func (t *TimeModify) SwitchDate(ctx context.Context, title string) {
	switch {
	case strings.Contains(title, "明天") || strings.Contains(title, "明早"):
		today := strings.Split(didaNow.String(), " ")[0]
		tomorrowTime := didaNow.AddDate(0, 0, 1)
		afterTime := didaNow.AddDate(0, 0, -1)
		tomorrow := strings.Split(tomorrowTime.String(), " ")[0]
		_, _, afd := afterTime.Date()

		go switchTime(ctx, title)
		select {
		case dayTime := <-stream:
			if strings.Contains(dayTime.(string), "T") {
				dStr := dayTime.(string)
				day := strings.Split(strings.Split(dStr, "T")[0], "-")
				dTime := strings.Split(dStr, "T")[0]
				dayNum, _ := strconv.Atoi(day[len(day)-1])
				if dayNum == afd {
					finalStream <- strings.Replace(dStr, dTime, today, -1)
				} else {
					tomTime := strings.Split(dayTime.(string), "T")[1]
					finalStream <- fmt.Sprintf("%sT%s", tomorrow, tomTime)
				}

			} else {
				finalStream <- fmt.Sprintf("%sT%s", tomorrow, dayTime)
			}

			return
		}

	case strings.Contains(title, "后天"):
		afterTime := didaNow.AddDate(0, 0, 2)
		after := strings.Split(afterTime.String(), " ")[0]
		go switchTime(ctx, title)
		select {
		case dayTime := <-stream:
			if strings.Contains(dayTime.(string), "T") {
				reTime := strings.Split(dayTime.(string), "T")[1]
				finalStream <- fmt.Sprintf("%sT%s", after, reTime)
			} else {
				finalStream <- fmt.Sprintf("%sT%s", after, dayTime)
			}
			return
		}
	case strings.Contains(title, "天后"):
		numStr := getStr(strings.Split(title, "天后")[0])
		num, _ := strconv.Atoi(numStr)
		date := didaNow.AddDate(0, 0, +num)
		after := strings.Split(date.String(), " ")[0]
		go switchTime(ctx, title)
		select {
		case dayTime := <-stream:
			if strings.Contains(dayTime.(string), "T") {
				reTime := strings.Split(dayTime.(string), "T")[1]
				finalStream <- fmt.Sprintf("%sT%s", after, reTime)
			} else {
				finalStream <- fmt.Sprintf("%sT%s", after, dayTime)
			}
			return
		}

	case strings.Contains(title, "月") && (strings.Contains(title, "号") || strings.Contains(title, "日")):
		var day string
		year := didaNow.Year()
		month := getStr(strings.Split(title, "月")[0])
		if strings.Contains(title, "号") {
			day = getStr(strings.Split(strings.Split(title, "号")[0], "月")[1])
		} else {
			day = getStr(strings.Split(strings.Split(title, "日")[0], "月")[1])
		}

		date := fmt.Sprintf("%v-%v-%v", year, month, day)
		go switchTime(ctx, title)
		select {
		case dayTime := <-stream:
			if strings.Contains(dayTime.(string), "T") {
				reTime := strings.Split(dayTime.(string), "T")[1]
				finalStream <- fmt.Sprintf("%sT%s", date, reTime)
			} else {
				finalStream <- fmt.Sprintf("%sT%s", date, dayTime)
			}
			return
		}

	default:
		go switchTime(ctx, title)
		select {
		case dayTime := <-stream:
			today := strings.Split(didaNow.String(), " ")[0]
			if strings.Contains(dayTime.(string), "T") {
				finalStream <- dayTime.(string)
			} else {

				finalStream <- fmt.Sprintf("%sT%s", today, dayTime.(string))
			}

			return
		}
	}

}

func switchTime(ctx context.Context, title string) {
	defer ctx.Done()
	switch {
	case (strings.Contains(title, "下午") || strings.Contains(title, "晚上")) && strings.Contains(title, "点半"):
		if strings.Contains(title, "晚上") {
			title = strings.Replace(title, "晚上", "午", -1)
		} else {
			title = title
		}
		timeStr := getStr(strings.Split(strings.Split(title, "点")[0], ("午"))[1])
		timeInt, _ := strconv.Atoi(timeStr)
		dTime := midday + timeInt
		stream <- fmt.Sprintf("%v:30:00%s", dTime, timeFlag)

	case (strings.Contains(title, "下午") || strings.Contains(title, "晚上")) && strings.Contains(title, "点") && !strings.Contains(title, "点半"):
		if strings.Contains(title, "晚上") {
			title = strings.Replace(title, "晚上", "午", -1)
		} else {
			title = title
		}
		timeStr := getStr(strings.Split(strings.Split(title, "点")[0], "午")[1])
		timeInt, _ := strconv.Atoi(timeStr)
		dTime := midday + timeInt

		stream <- fmt.Sprintf("%v:00:00%s", dTime, timeFlag)

	case strings.Contains(title, "晚上") && !strings.Contains(title, "点"):
		today := time.Now().Format("2006-01-02T15:04:05")
		nightTime := strings.Split(today, "T")[1]
		stream <- fmt.Sprintf("%v%s", strings.Replace(today, nightTime, defaultNight, -1), timeFlag)

	case strings.Contains(title, "分钟后"):
		var result string
		splitStr := strings.Split(title, "分钟后")[0]
		intStr := getSpecialStr(splitStr)
		strInt, _ := strconv.Atoi(intStr)
		result = didaNow.Add(time.Duration(strInt) * time.Minute).Format("2006-01-02T15:04:05")
		stream <- fmt.Sprintf("%v%s", result, timeFlag)

	case strings.Contains(title, "小时后"):
		var result string
		splitStr := strings.Split(title, "小时后")[0]
		intStr := getSpecialStr(splitStr)
		strInt, _ := strconv.Atoi(intStr)
		result = didaNow.Add(time.Duration(strInt) * time.Hour).Format("2006-01-02T15:04:05")
		stream <- fmt.Sprintf("%v%s", result, timeFlag)

	case strings.Contains(title, "：") || strings.Contains(title, ":"):
		var (
			hourStr string
			result  string
		)
		if strings.Contains(title, "：") {
			splitStr := strings.Split(title, "：")
			inputMin := reNum(splitStr[1])
			inputHourInt, _ := strconv.Atoi(reNum(splitStr[0]))
			y, m, d := didaNow.Date()
			if inputHourInt > 12 {
				hourInt := inputHourInt - 8
				result = fmt.Sprintf("%v-%v-%vT%v:%s:00%s", y, switchMonth(m.String()), d, hourInt, inputMin, timeFlag)

			} else {
				currentHour, _ := strconv.Atoi(strings.Split(strings.Split(didaNow.String(), " ")[1], ":")[0])
				if currentHour > 4 {
					hourStr = strconv.Itoa(inputHourInt + 4)
				} else {
					hourStr = strconv.Itoa(inputHourInt - 8)
				}
				result = fmt.Sprintf("%v-%v-%vT%s:%s:00%s", y, switchMonth(m.String()), d, hourStr, inputMin, timeFlag)
			}
			stream <- result

		} else {
			splitStr := strings.Split(title, ":")
			inputHour := reNum(splitStr[0])
			inputMin := reNum(splitStr[1])
			inputHourInt, _ := strconv.Atoi(inputHour)
			y, m, d := didaNow.Date()
			if inputHourInt > 12 {
				hourInt := inputHourInt - 8
				result = fmt.Sprintf("%v-%v-%vT%v:%s:00%s", y, switchMonth(m.String()), d, hourInt, inputMin, timeFlag)
			} else {
				currentHour, _ := strconv.Atoi(strings.Split(strings.Split(didaNow.String(), " ")[1], ":")[0])
				if currentHour > 4 {
					hourStr = strconv.Itoa(inputHourInt + 4)
				} else {
					hourStr = strconv.Itoa(inputHourInt - 8)
				}
				result = fmt.Sprintf("%v-%v-%vT%s:%s:00%s", y, switchMonth(m.String()), d, hourStr, inputMin, timeFlag)
			}
			stream <- result
		}

	case (strings.Contains(title, "点") && strings.Contains(title, "分")) && !strings.Contains(title, "点半") && !(strings.Contains(title, "下午") || strings.Contains(title, "晚上")):
		var hourStr string
		inputHour := getSpecialStr(strings.Split(title, "点")[0])
		inputMin := getSpecialStr(strings.Split(strings.Split(title, "点")[1], "分")[0])
		inputHourInt, _ := strconv.Atoi(inputHour)
		y, m, d := didaNow.Date()
		currentHour, _ := strconv.Atoi(strings.Split(strings.Split(didaNow.String(), " ")[1], ":")[0])

		if currentHour > 4 {
			hourStr = strconv.Itoa(inputHourInt + 4)
		} else {
			hourStr = inputHour
		}
		dateStr := fmt.Sprintf("%v-%v-%vT%s:%s:00", y, switchMonth(m.String()), d, hourStr, inputMin)
		stream <- fmt.Sprintf("%v%s", dateStr, timeFlag)

	case strings.Contains(title, "点") && !strings.Contains(title, "点半") && !(strings.Contains(title, "下午") || strings.Contains(title, "晚上")):
		var (
			//result  string
			timeStr string
		)
		splitStr := strings.Split(title, "点")
		inputHour := getStr(splitStr[0])
		inputHourInt, _ := strconv.Atoi(inputHour)

		y, m, d := didaNow.Date()
		currentHour, _ := strconv.Atoi(strings.Split(strings.Split(didaNow.String(), " ")[1], ":")[0])
		if currentHour > 4 {
			timeStr = strconv.Itoa(inputHourInt - 8)
		} else {
			timeStr = inputHour
		}
		dateStr := fmt.Sprintf("%v-%v-%vT%s:00:00", y, switchMonth(m.String()), d, timeStr)
		stream <- fmt.Sprintf("%v%s", dateStr, timeFlag)

	default:
		stream <- fmt.Sprintf("%s%s", didaNow.Format("2006-01-02T15:04:05"), timeFlag)
	}
}

func reNum(title string) (num string) {
	reg, _ := regexp.Compile("\\d+")
	num = reg.FindString(title)
	return
}

func getSpecialStr(origin string) (Fixed string) {
	length := len([]rune(origin))
	getNum := reNum(origin)
	if getNum != "" {
		return getNum
	}

	switch {
	case strings.Contains(origin, "十"):
		if length == 1 {
			Fixed = switchNum(origin)
			return
		} else if length == 2 {
			if judgeTwoNum(origin) {
				Fixed = switchNum(origin)
				return
			} else {
				step1 := strings.Split(origin, "十")[1]
				if judgeNum(step1) {
					Fixed = switchNum(origin)
					return
				}
			}
		} else if length >= 3 {
			sliceNum := origin[length-3:]
			splitNum := strings.Split(sliceNum, "十")
			if judgeNum(splitNum[0]) {
				Fixed = switchNum(sliceNum)
				return
			} else {
				Fixed = switchNum(sliceNum[1:])
				return
			}
		}

	default:
		num := origin[length-1:]
		Fixed = switchNum(num)
		return
	}

	return
}

func judgeNum(str string) bool {
	switch str {
	case "二":
		return true
	case "三":
		return true
	case "四":
		return true
	case "五":
		return true
	case "六":
		return true
	case "七":
		return true
	case "八":
		return true
	case "九":
		return true

	default:
		return false

	}
}

func judgeTwoNum(origin string) bool {
	switch origin {
	case "十一":
		return true
	case "十二":
		return true
	case "十三":
		return true
	case "十四":
		return true
	case "十五":
		return true
	case "十六":
		return true
	case "十七":
		return true
	case "十八":
		return true
	case "十九":
		return true
	default:
		return false
	}
}

func switchMonth(month string) string {
	switch month {
	case "January":
		return "1"
	case "February":
		return "2"
	case "March":
		return "3"
	case "April":
		return "4"
	case "May":
		return "5"
	case "June":
		return "6"
	case "July":
		return "7"
	case "August":
		return "8"
	case "September":
		return "9"
	case "October":
		return "10"
	case "November":
		return "11"
	case "December":
		return "12"
	default:
		return ""

	}
}

func getStr(origin string) (Fixed string) {
	numStr := switchNum(origin)
	if len(numStr) == 1 {
		Fixed = numStr[len(numStr)-1:]
		return
	} else if len(numStr) >= 2 {
		Fixed = numStr[len(numStr)-2:]
		return
	}
	return
}

func switchNum(title string) (num string) {
	lengthOfstr := []rune(title)
	getNum := reNum(title)
	if getNum != "" {
		return getNum
	}
	switch {
	case strings.Contains(title, "一") && len(lengthOfstr) == 1:
		return "1"
	case strings.Contains(title, "二") && len(lengthOfstr) == 1:
		return "2"
	case strings.Contains(title, "两") && len(lengthOfstr) == 1:
		return "2"
	case strings.Contains(title, "三") && len(lengthOfstr) == 1:
		return "3"
	case strings.Contains(title, "四") && len(lengthOfstr) == 1:
		return "4"
	case strings.Contains(title, "五") && len(lengthOfstr) == 1:
		return "5"
	case strings.Contains(title, "六") && len(lengthOfstr) == 1:
		return "6"
	case strings.Contains(title, "七") && len(lengthOfstr) == 1:
		return "7"
	case strings.Contains(title, "八") && len(lengthOfstr) == 1:
		return "8"
	case strings.Contains(title, "九") && len(lengthOfstr) == 1:
		return "9"
	case strings.Contains(title, "十") && len(lengthOfstr) == 1:
		return "10"
	case strings.Contains(title, "十一") && len(lengthOfstr) == 2:
		return "11"
	case strings.Contains(title, "十二") && len(lengthOfstr) == 2:
		return "12"
	case strings.Contains(title, "十三") && len(lengthOfstr) == 2:
		return "13"
	case strings.Contains(title, "十四") && len(lengthOfstr) == 2:
		return "14"
	case strings.Contains(title, "十五") && len(lengthOfstr) == 2:
		return "15"
	case strings.Contains(title, "十六") && len(lengthOfstr) == 2:
		return "16"
	case strings.Contains(title, "十七") && len(lengthOfstr) == 2:
		return "17"
	case strings.Contains(title, "十八") && len(lengthOfstr) == 2:
		return "18"
	case strings.Contains(title, "十九") && len(lengthOfstr) == 2:
		return "19"
	case strings.Contains(title, "二十") && len(lengthOfstr) == 2:
		return "20"
	case strings.Contains(title, "二十一") && len(lengthOfstr) == 3:
		return "21"
	case strings.Contains(title, "二十二") && len(lengthOfstr) == 3:
		return "22"
	case strings.Contains(title, "二十三") && len(lengthOfstr) == 3:
		return "23"
	case strings.Contains(title, "二十四") && len(lengthOfstr) == 3:
		return "24"
	case strings.Contains(title, "二十五") && len(lengthOfstr) == 3:
		return "25"
	case strings.Contains(title, "二十六") && len(lengthOfstr) == 3:
		return "26"
	case strings.Contains(title, "二十七") && len(lengthOfstr) == 3:
		return "27"
	case strings.Contains(title, "二十八") && len(lengthOfstr) == 3:
		return "28"
	case strings.Contains(title, "二十九") && len(lengthOfstr) == 3:
		return "29"
	case strings.Contains(title, "三十") && len(lengthOfstr) == 2:
		return "30"
	case strings.Contains(title, "三十一") && len(lengthOfstr) == 3:
		return "31"

	default:
		return title
	}

}

func (t *TimeModify) GetTime() string {
	defer close(stream)
	for {
		select {
		case data := <-finalStream:
			return data
		}
	}
}
