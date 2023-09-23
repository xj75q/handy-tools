## go-tools
### 【仓库作用】

本仓库主要存放一些go语言写的小工具，方便大家及自己使用



### 【使用方法】

#### 1> 文件重命名工具：
go renamefile --help 去查看里面有哪些参数需要输入

可用参数如下：
   --frompath, -f            Whether to take a value from the input file path,default 'false'  (default: false)
   --input value, -i value   please input file path
   --name value, -n value    please input the new name,if you take a value from path,then input the name key worlds (default: new)
   --output value, -o value  please input the output path
   --type value, -t value    please input the file type,eg: go, python ,txt (default: go)
   --help, -h                show help

eg:

```
  ./renamefile --input="/home/dir/go-dir" --frompath=true --name="day"
```



2> hideIcon 隐藏任务栏图标(只在win上可以使用)

使用方法：

```
alt+1 选择窗口；
alt+2 隐藏该程序的任务栏图标； 
alt+3 恢复任务栏图标
```

注：主要针对微信的闪闪闪，只是隐藏任务栏图标，在主窗口还能聊天

具体使用：

```
1> 微信点击置顶
2> 点击运行该软件hideIcon.exe
3> 鼠标指针移动到要选择的窗口中（比如微信）
4> 使用快捷键（其中任意2个）：
alt+1 选择窗口；
alt+2 隐藏该程序的任务栏图标； 
alt+3 恢复任务栏图标
5> 微信取消置顶
```

注：对于微信，需要先置顶再操作；对于其他软件，比如网易云，可以跳过1和5。（只在windows上有效果）

