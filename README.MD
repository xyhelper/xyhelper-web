# xyhelper-web

文档请查阅 [https://xyhelper.cn/web/](https://xyhelper.cn/web/)

## 环境变量

- PORT: 端口号，默认 8080
- GIN_MODE: gin 模式，默认 release, 可选 debug 或 release
- TIMEOUTMS: 超时时间(单位:毫秒)，默认 180000 毫秒
- AUTH_SECRET_KEY: 登陆密码 默认为空
- BASE_URI: 基础路径 默认为空   当指定时，会忽略前端传递的 BaseURI 参数
- ACCESS_TOKEN: 访问令牌 默认为空  当指定时，会忽略前端传递的 AccessToken 参数
- KFURL: 提示注册时显示的图片地址，默认为联系客服二维码地址