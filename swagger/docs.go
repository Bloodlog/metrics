// Package swagger Code generated by swaggo/swag. DO NOT EDIT
package swagger

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
        "/": {
            "get": {
                "description": "Генерирует HTML-страницу с перечнем метрик (gauge и counter)",
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "Info"
                ],
                "summary": "Список метрик",
                "responses": {
                    "200": {
                        "description": "HTML страница с метриками",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Info"
                ],
                "summary": "Проверка состояния сервиса",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/update": {
            "post": {
                "description": "Обновляет метрику с переданными параметрами",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Json"
                ],
                "summary": "Обновление метрики",
                "parameters": [
                    {
                        "description": "Metrics Update Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.MetricsUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Response with success status",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/update/{metricType}/{metricName}/{metricValue}": {
            "post": {
                "description": "Обновляет значение метрики по её типу (counter или gauge) на основе переданных параметров",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Text"
                ],
                "summary": "Обновление значения метрики",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Тип метрики (counter или gauge)",
                        "name": "metricType",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Имя метрики",
                        "name": "metricName",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Новое значение метрики",
                        "name": "metricValue",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Метрика успешно обновлена",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/updates": {
            "post": {
                "description": "Обновляет несколько метрик с переданными параметрами",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Json"
                ],
                "summary": "Обновление нескольких метрик",
                "parameters": [
                    {
                        "description": "Metrics Update Request List",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.MetricsUpdateRequest"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully updated",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/value": {
            "post": {
                "description": "Получает значение метрики по имени и типу",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Json"
                ],
                "summary": "Получение значения метрики",
                "parameters": [
                    {
                        "description": "Запрос на получение метрики",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.MetricsGetRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.MetricsResponse"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Метрика не найдена",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/value/{metricType}/{metricName}": {
            "get": {
                "description": "Возвращает значение метрики по ее типу (counter или gauge) в формате текста",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Text"
                ],
                "summary": "Получение значения метрики",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Тип метрики (counter или gauge)",
                        "name": "metricType",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Имя метрики",
                        "name": "metricName",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Метрика возвращена успешно",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Метрика не найдена",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.MetricsGetRequest": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "Имя метрики.",
                    "type": "string"
                },
                "type": {
                    "description": "Тип метрики: counter или gauge.",
                    "type": "string"
                }
            }
        },
        "dto.MetricsResponse": {
            "type": "object",
            "properties": {
                "delta": {
                    "description": "Значение counter.",
                    "type": "integer"
                },
                "id": {
                    "description": "Тип метрики: counter или gauge.",
                    "type": "string"
                },
                "type": {
                    "description": "Имя метрики.",
                    "type": "string"
                },
                "value": {
                    "description": "Значение gauge.",
                    "type": "number"
                }
            }
        },
        "dto.MetricsUpdateRequest": {
            "type": "object",
            "properties": {
                "delta": {
                    "description": "Значение counter.",
                    "type": "integer"
                },
                "id": {
                    "description": "Имя метрики.",
                    "type": "string"
                },
                "type": {
                    "description": "Тип метрики: counter или gauge.",
                    "type": "string"
                },
                "value": {
                    "description": "Значение gauge.",
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:8080",
	BasePath:         "/.",
	Schemes:          []string{},
	Title:            "Metrics API",
	Description:      "API для управления метриками",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
