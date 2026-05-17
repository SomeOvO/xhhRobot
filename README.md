# XhhRobot

小黑盒类Grok机器人

# 能做什么？

自动检查指定用户的@消息并使用Ai回复

- 自定义提示词

- 自定义Ai接口

# 开始使用

[手把手教程](https://blog.sakurasen.cn/post/1778819699353/)

## OpenAI 联网配置

如果使用 `gpt-5.4` 并希望模型联网，需要使用 OpenAI Responses API，并开启 `webSearch`。仅把模型名改为 `gpt-5.4` 不会自动联网。

```json
{
  "ai": {
    "model": "gpt-5.4",
    "baseUrl": "https://api.openai.com/v1/responses",
    "token": "你的 OpenAI API Key",
    "prompt": "你的提示词",
    "webSearch": true,
    "searchContextSize": "medium"
  }
}
```

如果必须继续使用 Chat Completions 接口，请将模型改为 `gpt-5-search-api`，并开启 `webSearch`。`gpt-5.4` 的联网能力推荐通过 Responses API 的 `web_search` 工具启用。

## 下载

前往[Release下载](https://github.com/SomeOvO/xhhRobot/releases)您对应的系统版本

# PR&Issues

欢迎各位提出Pr以及Issues。

在提pr前还是建议先去Issues请求一下，避免与其他人冲突。
