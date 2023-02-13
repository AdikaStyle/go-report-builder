# greypot

An experimental library for Report building and generations, based on [go-report-builder](https://github.com/AdikaStyle/go-report-builder).

Very much WIP.

## What does it do?

Say you want to produce reports or other such type of documents in your applications. 
`greypot` allows you to design your reports with HTML as template files that use the standard Go `html/template` templating engine.

These HTML reports can then be generated as HTML, PNG or PDF via endpoints that greypot adds to your application.

Once you add the middleware to your application, it adds the following routes:

```
GET /reports/list

GET /reports/preview/:reportTemplateName

GET /reports/render/:reportTemplateName

POST /reports/export/html/:reportTemplateName

POST /reports/export/png/:reportTemplateName

POST /reports/export/pdf/:reportTemplateName
```

You can then call these from within your applications to generate/export the reports e.g. from a frontend UI.


## Playwright Module

Currently, we are focusing on making the playwright based renderer work really good! The base project used Chrome Developer Protocol to connect with a Chromium instance. We are evaluating that decision [here](https://github.com/nndi-oss/greypot/issues/1)

In order to use the [Playwright](https://github.com/playwright-community/playwright-go) rendering functionality, you will need to have the [playwright dependencies](https://playwright.dev/docs/cli#install-system-dependencies) installed.

Read [here](https://playwright.dev/docs/cli#install-system-dependencies) for more info. But in short, you can use the following command to do so:

```
$ npx playwright install-deps chromium
```

### Basic Fiber example


```go
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

	module := greypot.NewModule(
		greypot.WithRenderTimeout(10*time.Second),
		greypot.WithViewport(2048, 1920),
		greypot.WithTemplateEngine(engine.NewDjangoTemplateEngine()),
		greypot.WithTemplatesRepository(greypot.NewFilesystemRepo("./templates/")),
		greypot.WithPlaywrightRenderer(),
	)

	greypotFiber.Use(app.Group("/api"), module)

	app.Listen(":3000")
}
```
