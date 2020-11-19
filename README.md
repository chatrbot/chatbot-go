## Chatbot 的 Go 版本 DEMO

### 接入流程和如何获取 token

点击[这里](https://github.com/chatrbot/chatbot)查看

### 运行示例

1. 下载代码  
   `git clone https://github.com/chatrbot/chatbot-go`
2. 安装 Go 依赖库
   `go mod download`
3. 替换代码中的`token`([如何获取我的 token](https://github.com/chatrbot/chatbot#faq))或者在命令行中添加 token 参数
   修改代码: `go run /example/main.go` 或  
   不修改代码: `go run /example/main.go -token your_token`
4. 发送“hello”给机器人,会返回一个“world”.

### 简单封装

里面包含了一个简易封装的 sdk,可直接使用
