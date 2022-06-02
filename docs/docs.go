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
        "/:deviceID/getList": {
            "post": {
                "description": "获取歌曲列表，0表示全部歌曲的列表，1表示已点歌曲的列表",
                "consumes": [
                    "application/json"
                ],
                "summary": "获取歌曲列表",
                "parameters": [
                    {
                        "description": "参数",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web_server.SongListRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web_server.SongListResponse"
                        }
                    }
                }
            }
        },
        "/:deviceID/orderSong": {
            "post": {
                "description": "点歌",
                "consumes": [
                    "application/json"
                ],
                "summary": "点歌",
                "parameters": [
                    {
                        "description": "参数",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web_server.OrderSongRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web_server.BaseResponse"
                        }
                    }
                }
            }
        },
        "/:deviceID/stickTopSong": {
            "post": {
                "description": "置顶已点歌曲",
                "consumes": [
                    "application/json"
                ],
                "summary": "置顶",
                "parameters": [
                    {
                        "description": "参数",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web_server.StickTopRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web_server.BaseResponse"
                        }
                    }
                }
            }
        },
        "/:deviceID/test": {
            "get": {
                "description": "测试",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "test"
                ],
                "summary": "测试",
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
        "web_server.BaseResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "web_server.OrderSongRequest": {
            "type": "object",
            "properties": {
                "song_id": {
                    "type": "integer"
                }
            }
        },
        "web_server.Song": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "singer_name": {
                    "type": "string"
                }
            }
        },
        "web_server.SongListRequest": {
            "type": "object",
            "properties": {
                "list_type": {
                    "type": "integer"
                }
            }
        },
        "web_server.SongListResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object",
                    "properties": {
                        "songs": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/web_server.Song"
                            }
                        }
                    }
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "web_server.StickTopRequest": {
            "type": "object",
            "properties": {
                "song_index": {
                    "type": "integer"
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