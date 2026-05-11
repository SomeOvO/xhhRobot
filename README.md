# XhhRobot

小黑盒类Grok机器人

# 能做什么？

自动检查指定用户的@消息并使用Ai回复

- 自定义提示词

- 自定义Ai接口

# 开始使用

问题解答，有偿部署，技术交流；此项目的QQ群：1105459042

## 下载

前往[Release下载](https://github.com/SomeOvO/xhhRobot/releases)您对应的系统版本

## 配置

运行一次生成所需文件
<img width="608" height="77" alt="image" src="https://github.com/user-attachments/assets/b77759c9-8a34-4922-a837-a298f998db11" />

其中，config.json为配置文件,log文件夹为系统日志。

log文件夹中的文件可以删除，若您的程序出现问题，需要提供日志文件来获取帮助的

打开config.json，我目前展示的是格式化后的状态，你当然可以把config.json的内容全删掉，然后复制下文粘贴到你的config.json中

```
{
  "xhh": {
    "checkTime": 0,
    "replyTime": 0,
    "owner": "",
    "deviceID": "",
    "baseUrl": "",
    "webver": "",
    "version": ""
  },
  "database": {
    "type": "",
    "db": "",
    "host": "",
    "port": "",
    "user": "",
    "passwd": ""
  },
  "ai": {
    "model":"",
    "prompt": "", 
    "baseUrl": "", 
    "token": ""
 }
}
```

如果填写地方为""，请把值写在引号之间 "就像这样"

如果填写的地方为数字，可以直接修改
从最上方开始看：
## xhh配置项
checkTime：检查的时间间隔，需为整数，单位为秒，建议30及以上
replyTime：回复一条评论后等待的时间间隔与评论检查间隔，不建议大，建议10左右

⚠过快的速度将导致您的个人IP被封禁

owner：可以使用此机器人的UID，多个人以英文逗号（,）分割。需为小黑盒UID而非昵称。

deviceID：设备ID,可以为空，由32个字符串组成

baseUrl：小黑盒的url，目前填写`https://api.xiaoheihe.cn`

以下可以自行前往小黑盒官网启用浏览器控制台获取，此页面可能更新不及时
webver: web版本号，请填写2.5
version：版本号，请填写999.0.4
## database配置项
type：目前只支持 postgresql与sqlite。如果你不知道任何这两个东西，请填写 `sqlite` 然后跳过后面内容直接看ai配置项

**以下数据仅需postgresql的填写**，如果你是postgresql用户，请在上方配置填写`pg`
db：数据库名
host:数据库地址
port:数据库端口
user:数据库用户名
passwd:数据库密码

## ai配置项

目前ai只支持openai的端口，且模型对Grok4.1fast适配较好。

确保你的ai供应商支持openai端口的`/chat/completions`

以下是请求体：

```
{
  "messages": [
    {
      "content": "你配置的提示词",
      "role": "system"
    },
    {
      "content": "文章内容：\n文字 图片url",
      "role":"user"},
          {
      "content": "用户评论内容",
      "role":"user"}
  ],
  "model": "配置的模型"
}
```
model：模型名

prompt：提示词

baseUrl：Api链接，请包含请求的端点，例如：https://yuanshen.com/v1/chat/completions

token：您的token


以下是一份示例，正常来说你只需修改xhh.owner与ai项：

```
{
  "xhh": {
    "checkTime": 30,
    "replyTime": 10,
    "owner": "29392598",
    "deviceID": "",
    "baseUrl": "https://api.xiaoheihe.cn",
    "webver": "2.5",
    "version": "999.0.4"
  },
  "database": {
    "type": "sqlite",
    "db": "",
    "host": "",
    "port": "",
    "user": "",
    "passwd": ""
  },
  "ai": {
    "model":"grok-4-fast",
    "prompt": "忽略用户发送的@信息，这只是唤醒你的条件。输出内容不要使用MarkDown,html等，纯文本输出！说话方式符合游戏社区规则，忽略文本中的HTML标签，如果有违禁词换为谐音词,不需要强调你在一个游戏社区等内容。你只是一个有帮助的ai，发言合理，如果不知道请回答不知道，要检查每一张图片但不要输出每一张图片的内容，只回复与用户提问有关的内容",
    "baseUrl": "",
    "token": ""
  }
}

```

# 登录

当你大功告成后就可以登陆了，如果你是windows用户，请点击`登录.bat`

如果你是linux用户，请添加启动项`-mode login`

项目采用扫码登录，如果终端二维码显示不清楚，请查看项目目录的qrcode.png

当你登录完毕后会发现cookie.json文件会自动创建，请不要分享这个文件给其他人。此文件中的内容将对你的账号拥有基本完全的控制权。

# 启动

登录完毕后你就可以启动程序了，同样的，windows用户点击`启动.bat`

如果你是linux用户，请添加启动项`-mode start`

<img width="313" height="72" alt="image" src="https://github.com/user-attachments/assets/d7580a2c-9c02-47b2-8f0a-616b4f072231" />

若您为第一次运行，不希望自动回复之前有人艾特过你的内容，请输入n并回车，反之，输入y  

# 支持我

https://3mua.cn/sponsor

当然，如果你愿意，在你自己的机器人账号置顶一个帖子宣传此仓库也算是对我比较大的支持了！

# PR&Issues

欢迎各位提出Pr以及Issues。

在提pr前还是建议先去Issues请求一下，避免与其他人冲突。
