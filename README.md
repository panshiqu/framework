# framework
a simple game framework

为方便只看到源码的朋友更好的理解这个简单游戏框架，请看我的博客：[中小型手机棋牌网络游戏服务端架构设计](https://blog.csdn.net/panshiqu/article/details/74572133)

已修改源码中Token，必要时可以联系我

## 修改

 + 数据库由golang database改[xorm](https://github.com/go-xorm/xorm)
 + @todo 增加redis缓存支持
 + 增加logrus日志组件
 
## 问题
每个game服务器只能对应一个类型和等级的游戏(如五子棋新手),这样有个问题就是如果某个游戏的阶段分布很少时会造成服务器资源的严重浪费。