package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "API для биржи раритетного мерча с аукционами, лотереями и конкурсами",
        "title": "Rare Merch Exchange API",
        "contact": {},
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/api/v1",
    "paths": {},
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Rare Merch Exchange API",
	Description:      "API для биржи раритетного мерча с аукционами, лотереями и конкурсами",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  doc,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
