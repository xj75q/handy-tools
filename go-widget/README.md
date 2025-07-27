
> go语言写的可能用到的一些小组件

包含：
```
封装log日志，使其不同等级带不同的颜色，方便排查问题使用
简单的不同限流实现方式
简单的redis分布式锁等
简单的基于etcd服务注册发现
```

##### 【colorLog】
使用方式:
```
c := colorLog.SetLog(colorLog.SetfileName("test1"), colorLog.SetInfoLoc(false))  
color := colorLog.NewLog(c)  
color.Errorf("这是error...")  
color.Infof("这里是info")  
color.Debugf("这里是debug")
```

##### 【etcd 注册发现】

在一个终端启动server
```
./server
```
另一个终端启动client
```
./client
```

