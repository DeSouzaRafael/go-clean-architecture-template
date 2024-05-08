// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/v0/user": {
            "post": {
                "description": "add new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Create User",
                "operationId": "create",
                "parameters": [
                    {
                        "description": "Set up users",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/input.UserInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/output.UserOutput"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/output.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/output.ResponseError"
                        }
                    }
                }
            }
        },
        "/v0/user/{id}": {
            "get": {
                "description": "Search for a user by its UUID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get User",
                "operationId": "getUser",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns the found user",
                        "schema": {
                            "$ref": "#/definitions/output.UserOutput"
                        }
                    },
                    "400": {
                        "description": "Invalid UUID format",
                        "schema": {
                            "$ref": "#/definitions/output.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/output.ResponseError"
                        }
                    }
                }
            },
            "put": {
                "description": "update existing user details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update User",
                "operationId": "updateUser",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update user details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/input.UserInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User Successfully updated"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/output.ResponseError"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/output.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/output.ResponseError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a user by its UUID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Delete User",
                "operationId": "deleteUser",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User successfully deleted"
                    },
                    "400": {
                        "description": "Invalid UUID format",
                        "schema": {
                            "$ref": "#/definitions/output.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/output.ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "input.UserInput": {
            "type": "object",
            "required": [
                "name",
                "phone"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "example": "user name"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "output.ResponseError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "message"
                }
            }
        },
        "output.UserOutput": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Go Clean Architecture Template API",
	Description:      "Template Golang",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
