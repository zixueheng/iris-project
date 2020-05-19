# Iris project
Iris + Gorm

# API

## 登录
请求：
URL：http://localhost:8080/adminapi/login
Method: Post
Body: `{"username":"admin","password":"123456"}`
响应：
```json
{
    "success": true,
    "code": 1000,
    "msg": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbl91c2VyX2lkIjoiMSIsImV4cCI6IjIwMjAtMDUtMTggMTU6Mzk6MDQiLCJpYXQiOiIyMDIwLTA1LTE4IDE1OjM2OjA0In0.8iMGRMMTR6j0WHV2xAZ6-qgNABGLa2SankV4iSDuo8A",
        "refresh_token": "pgdtfy56PtyRb2F59lwFbquK1Tnanz5EpfTwiBdPGD6BCEv2JHp8Kb6XJoQUpaE3"
    }
}
```

## 刷新Token
请求：
URL: http://localhost:8080/adminapi/refreshtoken
Method: Post
Header: authorization: bearer token(登录获取的token)
Body: `{"refresh_token": "pgdtfy56PtyRb2F59lwFbquK1Tnanz5EpfTwiBdPGD6BCEv2JHp8Kb6XJoQUpaE3"}`
响应：
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbl91c2VyX2lkIjoiMSIsImV4cCI6IjIwMjAtMDUtMTggMTU6NDY6NDYiLCJpYXQiOiIyMDIwLTA1LTE4IDE1OjQzOjQ2In0.Pfno5MrfeA0zKf0Db1qFiZI78Ir6KTzhSWxWFsf6DuQ",
        "refresh_token": "A61Wh4x9QuvSQAO8uxjHXAiCGK3eU9dRxrJKypkc7UtcCrFmnQRy8UNsMYSv1poz"
    }
}
```

