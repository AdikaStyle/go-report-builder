//go:build ignore
// +build ignore

package main

import (
	"embed"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot"
	greypotFiber "github.com/nndi-oss/greypot/http/fiber"
)

//go:embed "templates"
var templatesFS embed.FS

func main() {
	app := fiber.New()

	greypotModule := greypot.NewPlaywrightModule(10*time.Second, greypot.NewFSTemplateRepo(templatesFS))

	greypotFiber.Use(app, greypotModule)

	app.Listen(":3000")
}
