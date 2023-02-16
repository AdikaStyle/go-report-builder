package main

import (
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/nndi-oss/greypot"
	greypotFiber "github.com/nndi-oss/greypot/http/fiber"
	"github.com/nndi-oss/greypot/ui"
)

var templateDir string
var address string
var disableStudioUI bool = false

type UploadTemplateRequest struct {
	Name    string
	Content string
}

func init() {
	flag.StringVar(&templateDir, "templates", "./templates/", "Path to the directory with templates")
	flag.StringVar(&address, "listen", "localhost:7665", "Listen address for server")
	flag.BoolVar(&disableStudioUI, "disable-studio", false, "Disable the studio UI")
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

	if !disableStudioUI {
		studioTemplateStore := NewInmemoryTemplateRepository()
		studioModule := greypot.NewModule(
			greypot.WithRenderTimeout(10*time.Second),
			greypot.WithViewport(2048, 1920),
			greypot.WithDjangoTemplateEngine(),
			greypot.WithTemplatesRepository(studioTemplateStore),
			greypot.WithPlaywrightRenderer(),
		)

		studioRouter := app.Group("/_studio")

		greypotFiber.Use(studioRouter, studioModule)

		studioRouter.Post("/upload-template", func(c *fiber.Ctx) error {
			request := new(UploadTemplateRequest)
			err := c.BodyParser(&request)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message":    "failed to parse request body",
					"devMessage": err.Error(),
				})
			}
			nom := strings.TrimSpace(request.Name)
			err = studioTemplateStore.Add(nom, request.Content)

			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message":    "failed to upload template to store",
					"devMessage": err.Error(),
				})
			}

			return c.JSON(fiber.Map{
				"id":         nom,
				"message":    "uploaded the template successfully",
				"devMessage": "",
			})
		})

		frontendDistFS, err := fs.Sub(ui.FrontendFS, "dist")
		if err != nil {
			log.Fatalf("failed to read frontend assets dir got %v", err)
		}
		app.Use(filesystem.New(filesystem.Config{
			Root:   http.FS(frontendDistFS),
			Browse: false,
		}))
	}

	err = app.Listen(address)
	if err != nil {
		log.Fatalf("failed to start server at %s got %v", address, err)
	}
}
