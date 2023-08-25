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
