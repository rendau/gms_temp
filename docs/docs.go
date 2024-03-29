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
        "/config": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Get configs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.ConfigSt"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dopTypes.ErrRep"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Update configs",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/entities.ConfigSt"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dopTypes.ErrRep"
                        }
                    }
                }
            }
        },
        "/dic": {
            "get": {
                "description": "Get all dictionaries",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dic"
                ],
                "summary": "dictionaries",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.DicSt"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dopTypes.ErrRep"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dopTypes.ErrRep": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "string"
                },
                "error_code": {
                    "type": "string"
                },
                "fields": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                }
            }
        },
        "entities.ConfigContactsSt": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        },
        "entities.ConfigSt": {
            "type": "object",
            "properties": {
                "contacts": {
                    "$ref": "#/definitions/entities.ConfigContactsSt"
                }
            }
        },
        "entities.DicSt": {
            "type": "object"
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
