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
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/user/todo_categories": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves all user todo categories",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "operationId": "GetUserTodoCategories",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Creates a user todo category",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "operationId": "CreateUserTodoCategory",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/api/user/todo_categories/{id}": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves a user todo category by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "operationId": "GetUserTodoCategory",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Updates a user todo category with the specified id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "operationId": "UpdateUserTodoCategory",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Deletes a user todo category with the specified id",
                "tags": [
                    "Client"
                ],
                "operationId": "DeleteUserTodoCategory",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/version": {
            "get": {
                "description": "Gives the service version.",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Client"
                ],
                "operationId": "Version",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "AdminGroupAuth": {
            "type": "apiKey",
            "name": "GROUP",
            "in": "header"
        },
        "AdminUserAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header (add Bearer prefix to the Authorization value)"
        },
        "UserAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header (add Bearer prefix to the Authorization value)"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.1",
	Host:             "localhost",
	BasePath:         "/wellness",
	Schemes:          []string{"https"},
	Title:            "Rokwire Wellness Building Block API",
	Description:      "Rokwire Content Building Block API Documentation.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
