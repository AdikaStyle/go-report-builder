package handlers

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

func extractData(ctx *fiber.Ctx) (interface{}, error) {
	encodedData := ctx.Query("d")
	if encodedData == "" {
		return map[string]interface{}{}, nil
	}

	strData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to decode base64 data")
	}

	var jsonData interface{}
	if err := json.Unmarshal(strData, &jsonData); err != nil {
		return nil, stacktrace.Propagate(err, "failed to unmarshal data to json")
	}

	return jsonData, nil
}
