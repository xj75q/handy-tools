
### 【项目缘由】
由于wps和office 的很多功能都收费，所以就在前几年写了这些脚本，现在开源出来给需要的人使用
- pdf_combine.py 文件是合并pdf使用的，需要打印资料的时候可以用此合并
- ppt2word.py 文件是提取ppt当中文字的。当ppt中的内容需要整理到自己的电子笔记中，可使用此脚本
- full_ts2mp4.py 文件是将下载的ts视频整合成一个mp4文件的，可用于视频流的转换

### 【环境准备】
- mkvirtualenv py3office
- workon py3office
- pip install -r requirements.txt

### 【使用命令】
- python pdf_combine.py --help  可查看输入参数
- python pdf_combine.py -i "G:\\\小说pdf_转换前\\\" -o "G:\\\小说pdf_转换后\\\"  


**注：在win上使用，里面的路径为“\\\”(双斜杠)**


- send_task2dida.py脚本 (发送任务到滴答清单)
```
使用方法：
        1> 登录(先登录)
          python send_task2dida.py login <mail> <pwd>
        2> 配置默认项目（第二步设置默认项目）
          python send_task2dida.py project <projectName>
        3> 创建任务[只包含标题]（第三部创建任务）
          python send_task2dida.py <tasktitle>
        4> 创建send_task2dida.py 任务[标题+内容]
          python send_task2dida.py <tasktitle> <taskcontent>
        5> 创建任务[标题+剪切板内容]
```
所需其他包：
    - pip install tzlocal
    - pip install requests
    - pip install pyperclip

如果想在命令行使用此这个发送到滴答清单的python脚本，执行一下命令添加alias
```
vi ~/.bashrc
alias task="python /home/send2dida/send_task2dida.py" #添加此行
source ~/.bashrc
```

###  【其他】
使用full_ts2mp4.py 需在win上安装[ffmpeg](https://ffmpeg.org/) ，将py文件中的__init__里的ffmpeg路径修改掉


#### 后续会陆续更新其他脚本到此仓库

