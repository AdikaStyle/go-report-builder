//go:build ignore
// +build ignore

package main

import (
	"embed"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot"
	greypotFiber "github.com/nndi-oss/greypot/http/fiber"
	"github.com/nndi-oss/greypot/template/engine"
)

//go:embed "templates"
var templatesFS embed.FS

func main() {
	app := fiber.New()

	eng := engine.NewPongo2TemplateEngine()

	greypotModule := greypot.NewPlaywrightModuleWithCustomEngine(10*time.Second, greypot.NewFSTemplateRepo(templatesFS), eng)

	greypotFiber.Use(app, greypotModule)

	app.Listen(":3000")
}
