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
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/user-data": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Gets all related user data",
                "tags": [
                    "Client"
                ],
                "operationId": "GetUserData",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/UserDataResponse"
                        }
                    }
                }
            }
        },
        "/api/user/all_rings_records": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves all user ring record",
                "tags": [
                    "Client-RingsRecords"
                ],
                "operationId": "GetUserAllRingRecords",
                "parameters": [
                    {
                        "type": "string",
                        "description": "offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "limit - limit the result",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "order - Possible values: asc, desc. Default: desc",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "start_date - Start date filter in milliseconds as an integer epoch value",
                        "name": "start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "end_date - End date filter in milliseconds as an integer epoch value",
                        "name": "end_date",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/RingRecord"
                            }
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Deletes all user ring records (no matter of ring_id)",
                "tags": [
                    "Client-RingsRecords"
                ],
                "operationId": "DeleteAllUserRingRecords",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/user/rings": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves all user wellness ring entries",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client-Rings"
                ],
                "operationId": "GetUserRings",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/Ring"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Creates a user wellness ring entry",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client-Rings"
                ],
                "operationId": "CreateUserRing",
                "parameters": [
                    {
                        "description": "body json",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Ring"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Ring"
                        }
                    }
                }
            }
        },
        "/api/user/rings/{id}": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves a user wellness ring entry by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client-Rings"
                ],
                "operationId": "GetUserRing",
                "parameters": [
                    {
                        "description": "body json",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Ring"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Ring"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Deletes a user wellness ring entry with the specified id",
                "tags": [
                    "Client-Rings"
                ],
                "operationId": "DeleteUserRing",
                "responses": {}
            }
        },
        "/api/user/rings/{id}/history": {
            "post": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Creates a user wellness ring history entry",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client-Rings"
                ],
                "operationId": "CreateUserRingHistoryEntry",
                "parameters": [
                    {
                        "description": "body json",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/createUserRingRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/user/rings/{id}/history/{history-id}": {
            "delete": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Deletes a user wellness ring history entry with the specified id \u0026 history id",
                "tags": [
                    "Client-Rings"
                ],
                "operationId": "DeleteUserRingHistoryEntry",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/user/rings/{id}/records": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves all user ring record for a ring id",
                "tags": [
                    "Client-RingsRecords"
                ],
                "operationId": "GetUserRingRecords",
                "parameters": [
                    {
                        "type": "string",
                        "description": "offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "limit - limit the result",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "order - Possible values: asc, desc. Default: desc",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "start_date - Start date filter in milliseconds as an integer epoch value",
                        "name": "start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "end_date - End date filter in milliseconds as an integer epoch value",
                        "name": "end_date",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/RingRecord"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Creates a user ring record",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client-RingsRecords"
                ],
                "operationId": "CreateUserRingRecord",
                "parameters": [
                    {
                        "description": "body json",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/createUserRingRecordRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/RingRecord"
                            }
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Deletes all user ring record for a ring id",
                "tags": [
                    "Client-RingsRecords"
                ],
                "operationId": "DeleteUserRingRecords",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/user/rings/{id}/records/{record-id}": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves a user ring record by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client-RingsRecords"
                ],
                "operationId": "GetUserGetUserRingRecord",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/RingRecord"
                            }
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Updates a user ring record with the specified id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client-RingsRecords"
                ],
                "operationId": "UpdateUserRingRecord",
                "parameters": [
                    {
                        "description": "body json",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/RingRecord"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/RingRecord"
                            }
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Deletes a user ring record with the specified id",
                "tags": [
                    "Client-RingsRecords"
                ],
                "operationId": "DeleteUserRingRecord",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
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
                    "Client-TodoCategories"
                ],
                "operationId": "GetUserTodoCategories",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/TodoCategory"
                            }
                        }
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
                "tags": [
                    "Client-TodoCategories"
                ],
                "operationId": "CreateUserTodoCategory",
                "parameters": [
                    {
                        "description": "body json",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/TodoCategory"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/TodoCategory"
                        }
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
                    "Client-TodoCategories"
                ],
                "operationId": "GetUserTodoCategory",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/TodoCategory"
                        }
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
                    "Client-TodoCategories"
                ],
                "operationId": "UpdateUserTodoCategory",
                "parameters": [
                    {
                        "description": "body json",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/TodoCategory"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/TodoCategory"
                        }
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
                    "Client-TodoCategories"
                ],
                "operationId": "DeleteUserTodoCategory",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/user/todo_entries": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves all user todo entries",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client-TodoEntries"
                ],
                "operationId": "GetUserTodoEntries",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/TodoEntry"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Creates a user todo entry",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client-TodoEntries"
                ],
                "operationId": "CreateUserTodoEntry",
                "parameters": [
                    {
                        "description": "body json",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/TodoEntry"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/TodoEntry"
                        }
                    }
                }
            }
        },
        "/api/user/todo_entries/clear_completed_entries": {
            "delete": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Deletes all completed user todo entries",
                "tags": [
                    "Client-TodoEntries"
                ],
                "operationId": "DeleteCompletedUserTodoEntry",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/user/todo_entries/{id}": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves a user todo entry by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client-TodoEntries"
                ],
                "operationId": "GetUserTodoEntry",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/TodoEntry"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Updates a user todo entry with the specified id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client-TodoEntries"
                ],
                "operationId": "UpdateUserTodoEntry",
                "parameters": [
                    {
                        "description": "body json",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/TodoEntry"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/TodoEntry"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Deletes a user todo entry with the specified id",
                "tags": [
                    "Client-TodoEntries"
                ],
                "operationId": "DeleteUserTodoEntry",
                "responses": {
                    "200": {
                        "description": "OK"
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
                        "description": "OK"
                    }
                }
            }
        }
    },
    "definitions": {
        "CategoryRef": {
            "type": "object",
            "properties": {
                "app_id": {
                    "type": "string"
                },
                "color": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "org_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "Ring": {
            "type": "object",
            "properties": {
                "app_id": {
                    "type": "string"
                },
                "date_created": {
                    "type": "string"
                },
                "date_updated": {
                    "type": "string"
                },
                "history": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/RingHistoryEntry"
                    }
                },
                "id": {
                    "type": "string"
                },
                "org_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "RingHistoryEntry": {
            "type": "object",
            "properties": {
                "color_hex": {
                    "type": "string"
                },
                "date_created": {
                    "type": "string"
                },
                "date_updated": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "ring_id": {
                    "type": "string"
                },
                "unit": {
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "RingRecord": {
            "type": "object",
            "properties": {
                "app_id": {
                    "type": "string"
                },
                "date_created": {
                    "type": "string"
                },
                "date_updated": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "org_id": {
                    "type": "string"
                },
                "ring_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "RingRecordResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "RingResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "TodoCategory": {
            "type": "object",
            "properties": {
                "app_id": {
                    "type": "string"
                },
                "color": {
                    "type": "string"
                },
                "date_created": {
                    "type": "string"
                },
                "date_updated": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "org_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "TodoCategoryResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "TodoEntry": {
            "type": "object",
            "properties": {
                "app_id": {
                    "type": "string"
                },
                "category": {
                    "$ref": "#/definitions/CategoryRef"
                },
                "completed": {
                    "type": "boolean"
                },
                "date_created": {
                    "type": "string"
                },
                "date_updated": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "due_date_time": {
                    "type": "string"
                },
                "has_due_time": {
                    "type": "boolean"
                },
                "id": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "message_ids": {
                    "$ref": "#/definitions/model.MessageIDs"
                },
                "org_id": {
                    "type": "string"
                },
                "reminder_date_time": {
                    "type": "string"
                },
                "reminder_type": {
                    "type": "string"
                },
                "task_time": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                },
                "work_days": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "TodoEntryResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "UserDataResponse": {
            "type": "object",
            "properties": {
                "rings": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/RingResponse"
                    }
                },
                "rings_records": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/RingRecordResponse"
                    }
                },
                "todo_categories": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/TodoCategoryResponse"
                    }
                },
                "todo_entries": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/TodoEntryResponse"
                    }
                }
            }
        },
        "createUserRingRecordRequestBody": {
            "type": "object",
            "properties": {
                "ring_id": {
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "createUserRingRequestBody": {
            "type": "object",
            "properties": {
                "color_hex": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "unit": {
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "model.MessageIDs": {
            "type": "object",
            "properties": {
                "due_date_message_id": {
                    "type": "string"
                },
                "reminder_date_message_id": {
                    "type": "string"
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
        "InternalAPIAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header (add INTERNAL-API-KEY header with an appropriate value)"
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
	Version:          "1.0.2",
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
