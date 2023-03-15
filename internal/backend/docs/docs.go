// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/endpoints/add": {
            "post": {
                "description": "新增/修改路由配置（未携带id信息为新增）",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "endpoints"
                ],
                "summary": "新增/修改路由配置",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/endpoints/delete": {
            "post": {
                "description": "删除路由配置",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "endpoints"
                ],
                "summary": "删除路由配置",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "网关配置的id",
                        "name": "id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/endpoints/detail": {
            "post": {
                "description": "获取路由配置详情",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "endpoints"
                ],
                "summary": "获取路由配置详情",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "网关配置的id",
                        "name": "id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/endpoints/list": {
            "post": {
                "description": "获取所有路由配置",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "endpoints"
                ],
                "summary": "获取所有路由配置",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "页码数",
                        "name": "page_num",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 20,
                        "description": "分页尺寸",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/admin/login": {
            "post": {
                "description": "admin账号密码登录",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "admin账号密码登录",
                "operationId": "AdminLoginPassword",
                "parameters": [
                    {
                        "description": "auth info",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payload.AdminLoginPasswordReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/payload.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/payload.AdminLoginPasswordResp"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/auth/admin/register": {
            "post": {
                "description": "注册管理员用户",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "注册管理员用户",
                "parameters": [
                    {
                        "description": "auth info",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payload.AdminRegisterReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/payload.Response"
                        }
                    }
                }
            }
        },
        "/version": {
            "get": {
                "description": "获取当前版本",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "system"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "payload.AdminLoginPasswordReq": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "example": "123456"
                },
                "username": {
                    "type": "string",
                    "example": "admin"
                }
            }
        },
        "payload.AdminLoginPasswordResp": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "token_expire_at": {
                    "type": "integer"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "payload.AdminRegisterReq": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "example": "123456"
                },
                "username": {
                    "type": "string",
                    "example": "admin"
                }
            }
        },
        "payload.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "data": {},
                "msg": {
                    "type": "string"
                },
                "type": {
                    "description": "data msg error",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
