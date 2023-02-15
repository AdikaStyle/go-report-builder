package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot"
	greypotFiber "github.com/nndi-oss/greypot/http/fiber"
)

var templateDir string
var address string

func init() {
	flag.StringVar(&templateDir, "templates", "./templates/", "Path to the directory with templates")
	flag.StringVar(&address, "listen", "localhost:7665", "Listen address for server, defaults to 'localhost:7665'")
}

func main() {

	flag.Parse()

	if templateDir == "" {
		templateDir := os.Getenv("GREYPOT_TEMPLATE_DIR")
		if templateDir == "" {
			templateDir = "./templates/"
		}
	}

	absTemplateDir, err := filepath.Abs(templateDir)
	if err != nil {
		log.Fatalf("failed to get absolute path to templates, got %v", err)
	}

	entries, err := os.ReadDir(absTemplateDir)
	if err != nil {
		log.Fatalf("failed to read template directory got %v", err)
	}

	foundHTMLTemplates := false
	for _, e := range entries {
		if e.Type().IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".html") {
			foundHTMLTemplates = true
			break
		}
	}

	if !foundHTMLTemplates {
		log.Fatalf("Did not find any HTML template files in %s", templateDir)
	}

	if address == "" {
		address := os.Getenv("GREYPOT_ADDRESS")
		if address == "" {
			address = "localhost:7665"
		}
	}

	app := fiber.New()

	module := greypot.NewModule(
		greypot.WithRenderTimeout(10*time.Second),
		greypot.WithViewport(2048, 1920),
		greypot.WithDjangoTemplateEngine(),
		greypot.WithTemplatesFromFilesystem(absTemplateDir),
		greypot.WithPlaywrightRenderer(),
	)

	greypotFiber.Use(app, module)

	err = app.Listen(address)
	if err != nil {
		log.Fatalf("failed to start server at %s got %v", address, err)
	}
}
