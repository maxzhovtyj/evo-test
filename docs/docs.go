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
        "/api/load-data": {
            "post": {
                "description": "load data",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "api"
                ],
                "summary": "Load data to database from a csv file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "CSV file with data",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apperror.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apperror.Error"
                        }
                    }
                }
            }
        },
        "/api/transaction": {
            "get": {
                "description": "get transaction",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "api"
                ],
                "summary": "Gets a transaction",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Transaction id",
                        "name": "transactionId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Array with terminal ids, example: 3506,3507",
                        "name": "terminalIds",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Transaction status",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Transaction payment type",
                        "name": "paymentType",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Transaction min date post",
                        "name": "datePostFrom",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Transaction max date post",
                        "name": "datePostTo",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Transaction payment narrative",
                        "name": "paymentNarrative",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Transaction"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apperror.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apperror.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "apperror.Error": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.Transaction": {
            "type": "object",
            "properties": {
                "amountOriginal": {
                    "type": "number",
                    "example": 1
                },
                "amountTotal": {
                    "type": "number",
                    "example": 1
                },
                "commissionClient": {
                    "type": "number",
                    "example": 0
                },
                "commissionPS": {
                    "type": "number",
                    "example": 0
                },
                "commissionProvider": {
                    "type": "number",
                    "example": 0
                },
                "dateInput": {
                    "type": "string",
                    "example": "2022-08-12T11:25:27Z"
                },
                "datePost": {
                    "type": "string",
                    "example": "2022-08-12T14:25:27Z"
                },
                "partnerObjectId": {
                    "type": "integer",
                    "example": 1111
                },
                "payeeBankAccount": {
                    "type": "string",
                    "example": "UA713451373919523"
                },
                "payeeBankMFO": {
                    "type": "string",
                    "example": "254751"
                },
                "payeeId": {
                    "type": "integer",
                    "example": 14232155
                },
                "payeeName": {
                    "type": "string",
                    "example": "pumb"
                },
                "paymentNarrative": {
                    "type": "string",
                    "example": "?????????????????????????? ???????????? ???????????? ???????????????? ?????? ?????????????? ???????????? ??11/27122 ?????? 19.11.2020 ??."
                },
                "paymentNumber": {
                    "type": "string",
                    "example": "PS16698205"
                },
                "paymentType": {
                    "type": "string",
                    "example": "cash"
                },
                "requestId": {
                    "type": "integer",
                    "example": 20020
                },
                "service": {
                    "type": "string",
                    "example": "???????????????????? ????????????"
                },
                "serviceId": {
                    "type": "integer",
                    "example": 13980
                },
                "status": {
                    "type": "string",
                    "example": "accepted"
                },
                "terminalId": {
                    "type": "integer",
                    "example": 3506
                },
                "transactionId": {
                    "type": "integer",
                    "example": 1
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8089",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Evo Test Task",
	Description:      "API Server Evo Test Task",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
