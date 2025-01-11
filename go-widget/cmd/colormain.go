package main

import (
	"go-widget/colorLog"
)

func main() {
	//使用方式
	c := colorLog.SetLog(colorLog.SetfileName("test1"), colorLog.SetInfoLoc(false))
	color := colorLog.NewLog(c)
	color.Errorf("这是error...")
	color.Infof("这里是info")
	color.Debugf("这里是debug")
}
