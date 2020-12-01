## ChatBot 的 Go 版本 DEMO

### 接入流程和如何获取 token

点击[这里](https://github.com/chatrbot/chatbot)查看

### 运行

1. 下载代码  
   `git clone https://github.com/chatrbot/chatbot-go`
2. 安装 Go 依赖库
   `go mod download`  
3. [注册](https://www.ownthink.com/) 获取的一个智能机器人API的Token
4. 替换代码中的`token`([如何获取我的 token](https://github.com/chatrbot/chatbot#faq))或者在命令行中添加 token 参数
   修改代码: `go run /example/ai/main.go` 或  
   不修改代码: `go run /example/ai/main.go -token your_token -ai ai_token`
5. 就能实现一个能够只能回复的机器人

### 示例插件
1. AI机器人  
该示例需要自行注册思知的AI API Token，可以在这个[地址](https://www.ownthink.com/)获取  
下载代码并且成功[获取到ChatBot Token](https://github.com/chatrbot/chatbot#faq)  
    ```shell
    go run example/ai/main.go -token {ChatBotToken} -ai {AIToken}
    ```
    这时候和机器人对话，即可收到智能回复（如果是群内，需要@机器人）
    
2. 复读机
复读机插件并没有实际作用。只是作为一个SDK能力的展示和实际代码示例。展示了文本、图片、视频、语音类型的消息收发能力。
