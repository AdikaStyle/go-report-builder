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

	module := greypot.NewModule(
		greypot.WithRenderTimeout(10*time.Second),
		greypot.WithViewport(2048, 1920),
		greypot.WithTemplateEngine(engine.NewDjangoTemplateEngine()),
		greypot.WithTemplatesFromFilesytem("./templates/"),
		greypot.WithPlaywrightRenderer(),
	)

	greypotFiber.Use(app, module)

	embedModule := greypot.NewModule(
		greypot.WithRenderTimeout(10*time.Second),
		greypot.WithViewport(2048, 1920),
		greypot.WithTemplateEngine(engine.NewGolangTemplateEngine()),
		greypot.WithTemplatesFromFS(templatesFS),
		greypot.WithPlaywrightRenderer(),
	)

	greypotFiber.Use(app.Group("/embedded/"), embedModule)

	app.Listen(":3000")
}
