# Iris project 后台登录菜单权限角色管理
Iris MVC 多模块设计 + Gorm + MySQL + Redis

# API

## 公共

### 登录
- Request:
- URL：http://localhost:8080/adminapi/login
- Method: Post
- Body: `{"username":"admin","password":"123456"}`
- Response:
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

### 刷新Token
- Request:
- URL: http://localhost:8080/adminapi/refreshtoken
- Method: Post
- Header: authorization: bearer token(登录获取的token)
- Body: `{"refresh_token": "pgdtfy56PtyRb2F59lwFbquK1Tnanz5EpfTwiBdPGD6BCEv2JHp8Kb6XJoQUpaE3"}`
- Response:
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

## 管理员

### 管理员列表
- Request:
- URL: http://localhost:8080/adminapi/adminuser/list/1/2?username=admin (1是第1页，2是每页2条，username根据用户名查询)
- Method: Get
- Header: authorization: bearer token(登录获取的token)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": {
        "list": [
            {
                "id": 1,
                "created_at": "2020-05-23 09:29:08",
                "updated_at": "2020-05-23 09:29:08",
                "username": "admin",
                "role_id": 1,
                "role": {
                    "id": 1,
                    "created_at": "2020-05-23 09:29:08",
                    "name": "超级管理员",
                    "tag": "superadmin",
                    "status": 1
                },
                "phone": "15215657185",
                "status": 1
            },
            {
                "id": 2,
                "created_at": "2020-05-23 09:29:09",
                "updated_at": "2020-05-23 09:29:09",
                "username": "subadmin",
                "role_id": 2,
                "role": {
                    "id": 2,
                    "created_at": "2020-05-23 09:29:09",
                    "name": "子管理员",
                    "tag": "goods_manager",
                    "status": 1
                },
                "phone": "13721047437",
                "status": 1
            }
        ],
        "total": 2
    }
}
```


### 管理员详情
- Request:
- URL: http://localhost:8080/adminapi/adminuser/2 (2是id，如果是超级管理员不加载菜单树，超级管理员拥有所有菜单权限)
- Method: Get
- Header: authorization: bearer token(登录获取的token)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": {
        "id": 2,
        "created_at": "2020-05-23 09:29:09",
        "updated_at": "2020-05-23 09:29:09",
        "username": "subadmin",
        "role_id": 2,
        "role": {
            "id": 2,
            "created_at": "2020-05-23 09:29:09",
            "name": "子管理员",
            "tag": "goods_manager",
            "menus_tree": [
                {
                    "id": 1,
                    "created_at": "2020-05-23 09:29:09",
                    "p_id": 0,
                    "name": "主页",
                    "icon": "md-home",
                    "type": "menu",
                    "menu_path": "/admin/home/",
                    "api_path": "",
                    "method": "",
                    "sort": 1,
                    "status": 1,
                    "children": [
                        {
                            "id": 2,
                            "created_at": "2020-05-23 09:29:09",
                            "p_id": 1,
                            "name": "首页统计接口",
                            "icon": "",
                            "type": "api",
                            "menu_path": "",
                            "api_path": "/adminapi/statistic",
                            "method": "GET",
                            "sort": 0,
                            "status": 1,
                            "children": []
                        }
                    ]
                },
                {
                    "id": 3,
                    "created_at": "2020-05-23 09:29:09",
                    "p_id": 0,
                    "name": "权限管理",
                    "icon": "md-settings",
                    "type": "menu",
                    "menu_path": "/admin/setting/",
                    "api_path": "",
                    "method": "",
                    "sort": 2,
                    "status": 1,
                    "children": [
                        {
                            "id": 4,
                            "created_at": "2020-05-23 09:29:09",
                            "p_id": 3,
                            "name": "管理员",
                            "icon": "",
                            "type": "menu",
                            "menu_path": "/admin/setting/admin_user",
                            "api_path": "",
                            "method": "",
                            "sort": 0,
                            "status": 1,
                            "children": [
                                {
                                    "id": 5,
                                    "created_at": "2020-05-23 09:29:09",
                                    "p_id": 4,
                                    "name": "管理员列表接口",
                                    "icon": "",
                                    "type": "api",
                                    "menu_path": "",
                                    "api_path": "/adminapi/adminuser/list/%v/%v",
                                    "method": "GET",
                                    "sort": 0,
                                    "status": 1,
                                    "children": []
                                },
                                {
                                    "id": 6,
                                    "created_at": "2020-05-23 09:29:09",
                                    "p_id": 4,
                                    "name": "管理员详情接口",
                                    "icon": "",
                                    "type": "api",
                                    "menu_path": "",
                                    "api_path": "/adminapi/adminuser/%v",
                                    "method": "GET",
                                    "sort": 0,
                                    "status": 1,
                                    "children": []
                                },
                                {
                                    "id": 7,
                                    "created_at": "2020-05-23 09:29:09",
                                    "p_id": 4,
                                    "name": "管理员添加编辑接口",
                                    "icon": "",
                                    "type": "api",
                                    "menu_path": "",
                                    "api_path": "/adminapi/adminuser",
                                    "method": "POST",
                                    "sort": 0,
                                    "status": 1,
                                    "children": []
                                },
                                {
                                    "id": 8,
                                    "created_at": "2020-05-23 09:29:09",
                                    "p_id": 4,
                                    "name": "管理员删除接口",
                                    "icon": "",
                                    "type": "api",
                                    "menu_path": "",
                                    "api_path": "/adminapi/adminuser/%v",
                                    "method": "DELETE",
                                    "sort": 0,
                                    "status": 1,
                                    "children": []
                                },
                                {
                                    "id": 9,
                                    "created_at": "2020-05-23 09:29:09",
                                    "p_id": 4,
                                    "name": "管理员禁用启用接口",
                                    "icon": "",
                                    "type": "api",
                                    "menu_path": "",
                                    "api_path": "/adminapi/adminuser/status/%v",
                                    "method": "GET",
                                    "sort": 0,
                                    "status": 1,
                                    "children": []
                                }
                            ]
                        },
                        {
                            "id": 10,
                            "created_at": "2020-05-23 09:29:09",
                            "p_id": 3,
                            "name": "角色",
                            "icon": "",
                            "type": "menu",
                            "menu_path": "/admin/setting/role",
                            "api_path": "",
                            "method": "",
                            "sort": 0,
                            "status": 1,
                            "children": [
                                {
                                    "id": 11,
                                    "created_at": "2020-05-23 09:29:09",
                                    "p_id": 10,
                                    "name": "角色列表接口",
                                    "icon": "",
                                    "type": "api",
                                    "menu_path": "",
                                    "api_path": "/adminapi/role/list/%v/%v",
                                    "method": "GET",
                                    "sort": 0,
                                    "status": 1,
                                    "children": []
                                },
                                {
                                    "id": 12,
                                    "created_at": "2020-05-23 09:29:09",
                                    "p_id": 10,
                                    "name": "角色添加编辑接口",
                                    "icon": "",
                                    "type": "api",
                                    "menu_path": "",
                                    "api_path": "/adminapi/role",
                                    "method": "POST",
                                    "sort": 0,
                                    "status": 1,
                                    "children": []
                                }
                            ]
                        }
                    ]
                }
            ],
            "status": 1
        },
        "phone": "13721047437",
        "status": 1
    }
}
```

### 新增或编辑管理员
- Request:
- URL: http://localhost:8080/adminapi/adminuser
- Method: Post
- Header: authorization: bearer token(登录获取的token)
- Body: `{"id":14,"username":"editor","password":"7777777","phone":"16337734722","role_id":2,"status":1}` (json里面有`id`字段就是编辑、没有就是新增)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": null
}
```

### 删除管理员
- Request:
- URL: http://localhost:8080/adminapi/adminuser/14 (14是id)
- Method: Delete
- Header: authorization: bearer token(登录获取的token)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": null
}
```

### 禁用或启用管理员
- Request:
- URL: http://localhost:8080/adminapi/adminuser/status/2 (2是id，访问一次状态就修改成相反的)
- Method: Get
- Header: authorization: bearer token(登录获取的token)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": null
}
```

## 菜单

### 菜单树
- Request:
- URL: http://localhost:8080/adminapi/menu/tree
- Method: Get
- Header: authorization: bearer token(登录获取的token)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": [
        {
            "id": 1,
            "created_at": "2020-05-23 09:29:09",
            "p_id": 0,
            "name": "主页",
            "icon": "md-home",
            "type": "menu",
            "menu_path": "/admin/home/",
            "api_path": "",
            "method": "",
            "sort": 1,
            "status": 1,
            "children": [
                {
                    "id": 2,
                    "created_at": "2020-05-23 09:29:09",
                    "p_id": 1,
                    "name": "首页统计接口",
                    "icon": "",
                    "type": "api",
                    "menu_path": "",
                    "api_path": "/adminapi/statistic",
                    "method": "GET",
                    "sort": 0,
                    "status": 1,
                    "children": []
                }
            ]
        },
        {
            "id": 3,
            "created_at": "2020-05-23 09:29:09",
            "p_id": 0,
            "name": "权限管理",
            "icon": "md-settings",
            "type": "menu",
            "menu_path": "/admin/setting/",
            "api_path": "",
            "method": "",
            "sort": 2,
            "status": 1,
            "children": [
                {
                    "id": 4,
                    "created_at": "2020-05-23 09:29:09",
                    "p_id": 3,
                    "name": "管理员",
                    "icon": "",
                    "type": "menu",
                    "menu_path": "/admin/setting/admin_user",
                    "api_path": "",
                    "method": "",
                    "sort": 0,
                    "status": 1,
                    "children": [
                        {
                            "id": 5,
                            "created_at": "2020-05-23 09:29:09",
                            "p_id": 4,
                            "name": "管理员列表接口",
                            "icon": "",
                            "type": "api",
                            "menu_path": "",
                            "api_path": "/adminapi/adminuser/list/%v/%v",
                            "method": "GET",
                            "sort": 0,
                            "status": 1,
                            "children": []
                        },
                        {
                            "id": 6,
                            "created_at": "2020-05-23 09:29:09",
                            "p_id": 4,
                            "name": "管理员详情接口",
                            "icon": "",
                            "type": "api",
                            "menu_path": "",
                            "api_path": "/adminapi/adminuser/%v",
                            "method": "GET",
                            "sort": 0,
                            "status": 1,
                            "children": []
                        },
                        {
                            "id": 7,
                            "created_at": "2020-05-23 09:29:09",
                            "p_id": 4,
                            "name": "管理员添加编辑接口",
                            "icon": "",
                            "type": "api",
                            "menu_path": "",
                            "api_path": "/adminapi/adminuser",
                            "method": "POST",
                            "sort": 0,
                            "status": 1,
                            "children": []
                        },
                        {
                            "id": 8,
                            "created_at": "2020-05-23 09:29:09",
                            "p_id": 4,
                            "name": "管理员删除接口",
                            "icon": "",
                            "type": "api",
                            "menu_path": "",
                            "api_path": "/adminapi/adminuser/%v",
                            "method": "DELETE",
                            "sort": 0,
                            "status": 1,
                            "children": []
                        },
                        {
                            "id": 9,
                            "created_at": "2020-05-23 09:29:09",
                            "p_id": 4,
                            "name": "管理员禁用启用接口",
                            "icon": "",
                            "type": "api",
                            "menu_path": "",
                            "api_path": "/adminapi/adminuser/status/%v",
                            "method": "GET",
                            "sort": 0,
                            "status": 1,
                            "children": []
                        }
                    ]
                },
                {
                    "id": 10,
                    "created_at": "2020-05-23 09:29:09",
                    "p_id": 3,
                    "name": "角色",
                    "icon": "",
                    "type": "menu",
                    "menu_path": "/admin/setting/role",
                    "api_path": "",
                    "method": "",
                    "sort": 0,
                    "status": 1,
                    "children": [
                        {
                            "id": 11,
                            "created_at": "2020-05-23 09:29:09",
                            "p_id": 10,
                            "name": "角色列表接口",
                            "icon": "",
                            "type": "api",
                            "menu_path": "",
                            "api_path": "/adminapi/role/list/%v/%v",
                            "method": "GET",
                            "sort": 0,
                            "status": 1,
                            "children": []
                        },
                        {
                            "id": 12,
                            "created_at": "2020-05-23 09:29:09",
                            "p_id": 10,
                            "name": "角色添加编辑接口",
                            "icon": "",
                            "type": "api",
                            "menu_path": "",
                            "api_path": "/adminapi/role",
                            "method": "POST",
                            "sort": 0,
                            "status": 1,
                            "children": []
                        }
                    ]
                }
            ]
        }
    ]
}
```

### 新增或更新菜单
- Request:
- URL: http://localhost:8080/adminapi/menu
- Method: Post
- Header: authorization: bearer token(登录获取的token)
- Body: `{"id":12,"p_id": 10,"name": "角色添加编辑接口","icon": "","type": "api","menu_path": "","api_path": "/adminapi/role","method": "POST","sort": 0,"status": 1}` (json里面有`id`字段就是编辑、没有就是新增)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": null
}
```

### 删除菜单（包含所有子菜单）
- Request:
- URL: http://localhost:8080/adminapi/menu/10 (10是要删除的菜单ID)
- Method: Delete
- Header: authorization: bearer token(登录获取的token)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": null
}
```

### 禁用或启用菜单
- Request:
- URL: http://localhost:8080/adminapi/menu/status/3 (3是菜单ID)
- Method: Get
- Header: authorization: bearer token(登录获取的token)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": null
}
```

## 角色

### 角色列表
- Request:
- URL: http://localhost:8080/adminapi/role/list/1/10?name=超级管理员
- Method: Get
- Header: authorization: bearer token(登录获取的token)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": {
        "list": [
            {
                "id": 1,
                "created_at": "2020-05-23 09:29:08",
                "name": "超级管理员",
                "tag": "superadmin",
                "status": 1
            },
            {
                "id": 2,
                "created_at": "2020-05-23 09:29:09",
                "name": "子管理员",
                "tag": "goods_manager",
                "status": 1
            }
        ],
        "total": 2
    }
}
```

### 角色详情
- Request:
- URL: http://localhost:8080/adminapi/role/2 (超级管理员角色不加载菜单树)
- Method: Get
- Header: authorization: bearer token(登录获取的token)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": {
        "id": 2,
        "created_at": "2020-05-23 09:29:09",
        "name": "子管理员",
        "tag": "goods_manager",
        "menus_tree": [
            {
                "id": 1,
                "created_at": "2020-05-23 09:29:09",
                "p_id": 0,
                "name": "主页",
                "icon": "md-home",
                "type": "menu",
                "menu_path": "/admin/home/",
                "api_path": "",
                "method": "",
                "sort": 1,
                "status": 1,
                "children": [
                    {
                        "id": 2,
                        "created_at": "2020-05-23 09:29:09",
                        "p_id": 1,
                        "name": "首页统计接口",
                        "icon": "",
                        "type": "api",
                        "menu_path": "",
                        "api_path": "/adminapi/statistic",
                        "method": "GET",
                        "sort": 0,
                        "status": 1,
                        "children": []
                    }
                ]
            },
            {
                "id": 3,
                "created_at": "2020-05-23 09:29:09",
                "p_id": 0,
                "name": "权限管理",
                "icon": "md-settings",
                "type": "menu",
                "menu_path": "/admin/setting/",
                "api_path": "",
                "method": "",
                "sort": 2,
                "status": 1,
                "children": [
                    {
                        "id": 4,
                        "created_at": "2020-05-23 09:29:09",
                        "p_id": 3,
                        "name": "管理员",
                        "icon": "",
                        "type": "menu",
                        "menu_path": "/admin/setting/admin_user",
                        "api_path": "",
                        "method": "",
                        "sort": 0,
                        "status": 1,
                        "children": [
                            {
                                "id": 5,
                                "created_at": "2020-05-23 09:29:09",
                                "p_id": 4,
                                "name": "管理员列表接口",
                                "icon": "",
                                "type": "api",
                                "menu_path": "",
                                "api_path": "/adminapi/adminuser/list/%v/%v",
                                "method": "GET",
                                "sort": 0,
                                "status": 1,
                                "children": []
                            },
                            {
                                "id": 6,
                                "created_at": "2020-05-23 09:29:09",
                                "p_id": 4,
                                "name": "管理员详情接口",
                                "icon": "",
                                "type": "api",
                                "menu_path": "",
                                "api_path": "/adminapi/adminuser/%v",
                                "method": "GET",
                                "sort": 0,
                                "status": 1,
                                "children": []
                            },
                            {
                                "id": 7,
                                "created_at": "2020-05-23 09:29:09",
                                "p_id": 4,
                                "name": "管理员添加编辑接口",
                                "icon": "",
                                "type": "api",
                                "menu_path": "",
                                "api_path": "/adminapi/adminuser",
                                "method": "POST",
                                "sort": 0,
                                "status": 1,
                                "children": []
                            },
                            {
                                "id": 8,
                                "created_at": "2020-05-23 09:29:09",
                                "p_id": 4,
                                "name": "管理员删除接口",
                                "icon": "",
                                "type": "api",
                                "menu_path": "",
                                "api_path": "/adminapi/adminuser/%v",
                                "method": "DELETE",
                                "sort": 0,
                                "status": 1,
                                "children": []
                            },
                            {
                                "id": 9,
                                "created_at": "2020-05-23 09:29:09",
                                "p_id": 4,
                                "name": "管理员禁用启用接口",
                                "icon": "",
                                "type": "api",
                                "menu_path": "",
                                "api_path": "/adminapi/adminuser/status/%v",
                                "method": "GET",
                                "sort": 0,
                                "status": 1,
                                "children": []
                            }
                        ]
                    },
                    {
                        "id": 10,
                        "created_at": "2020-05-23 09:29:09",
                        "p_id": 3,
                        "name": "角色",
                        "icon": "",
                        "type": "menu",
                        "menu_path": "/admin/setting/role",
                        "api_path": "",
                        "method": "",
                        "sort": 0,
                        "status": 1,
                        "children": [
                            {
                                "id": 11,
                                "created_at": "2020-05-23 09:29:09",
                                "p_id": 10,
                                "name": "角色列表接口",
                                "icon": "",
                                "type": "api",
                                "menu_path": "",
                                "api_path": "/adminapi/role/list/%v/%v",
                                "method": "GET",
                                "sort": 0,
                                "status": 1,
                                "children": []
                            },
                            {
                                "id": 12,
                                "created_at": "2020-05-23 09:29:09",
                                "p_id": 10,
                                "name": "角色添加编辑接口",
                                "icon": "",
                                "type": "api",
                                "menu_path": "",
                                "api_path": "/adminapi/role",
                                "method": "POST",
                                "sort": 0,
                                "status": 1,
                                "children": []
                            }
                        ]
                    }
                ]
            }
        ],
        "status": 1
    }
}
```


### 角色创建或更新
- Request:
- URL: http://localhost:8080/adminapi/role
- Method: Post
- Header: authorization: bearer token(登录获取的token)
- Body: `{"id":3,"name":"信息管理员","tag":"editor","status":1,"menu_ids":[3,10,11,12]}` (json里面有`id`字段就是编辑、没有就是新增，menu_ids要确保从1级菜单往下衍生)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": null
}
```

### 角色删除
- Request:
- URL: http://localhost:8080/adminapi/role/3
- Method: Delete
- Header: authorization: bearer token(登录获取的token)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": null
}
```

### 角色禁用或启用
- Request:
- URL: http://localhost:8080/adminapi/role/status/2
- Method: Get
- Header: authorization: bearer token(登录获取的token)
- Response:
```json
{
    "success": true,
    "code": 200,
    "msg": "成功",
    "data": null
}
```