package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/nndi-oss/greypot"
	greypotFiber "github.com/nndi-oss/greypot/http/fiber"
	"github.com/nndi-oss/greypot/ui"
)

var templateDir string
var host string
var port int
var disableStudioUI bool = false

type UploadTemplateRequest struct {
	Name    string
	Content string
}

func init() {
	flag.StringVar(&templateDir, "templates", "./templates/", "Path to the directory with templates")
	flag.StringVar(&host, "host", "0.0.0.0", "Host for server")
	flag.IntVar(&port, "port", 0, "Port for server (defaults to reading from $PORT or fallback to 7665)")
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

	if host == "" {
		address := os.Getenv("GREYPOT_HOST")
		if address == "" {
			address = "0.0.0.0"
		}
	}

	if port == 0 {
		envPort := os.Getenv("PORT")
		if envPort == "" {
			port = 7665
		} else {
			if parsedPort, err := strconv.Atoi(envPort); err != nil {
				port = 7665
			} else {
				port = parsedPort
			}
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

	err = app.Listen(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("failed to start server at %s:%d got %v", host, port, err)
	}
}
