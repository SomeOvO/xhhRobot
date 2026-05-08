# XhhRobot

目前暂无用处。

当前功能：

自动检查艾特列表（暂无分页）每十秒检查最新的20个@。

每五秒回复1个@。

# 配置文件

```json
{
  "xhh": {
    "baseUrl": "https://api.xiaoheihe.cn",
    "webver": "2.5",
    "version": "999.0.4"
  },
  "database": {
    "db": "数据库名",
    "host": "数据库地址",
    "port": "端口",
    "user": "用户名",
    "passwd": "密码"
  }
}
```

> xhh 中的 baseUrl若无特殊配置，均使用 **https://api.xiaoheihe.cn**
> webver与version请前往官网查看



# 项目结构

```
C:.
├─config -> 配置文件模块
├─log -> 日志目录（启动会自动生成）
├─loger -> 日志模块
├─pg ->数据库模块
└─xhh -> 小黑盒模块
```

