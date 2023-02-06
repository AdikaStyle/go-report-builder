package playwright

import (
	"fmt"
	"time"

	"github.com/nndi-oss/greypot/models"
	"github.com/playwright-community/playwright-go"
)

type pngReportExporter struct {
	timeout time.Duration
	// viewportHeight int
	// viewportWidth  int
}

func NewPngReportExporter(timeout time.Duration) *pngReportExporter {
	return &pngReportExporter{timeout: timeout}
}

func (pre *pngReportExporter) Export(url string, renderedTemplate []byte) ([]byte, *models.PrintOptions, error) {

	pw, err := playwright.Run()
	if err != nil {
		return nil, nil, fmt.Errorf("could not launch playwright: %w", err)
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		return nil, nil, fmt.Errorf("could not launch Chromium: %w", err)
	}

	context, err := browser.NewContext()
	if err != nil {
		return nil, nil, fmt.Errorf("could not create context: %w", err)
	}
	page, err := context.NewPage()
	if err != nil {
		return nil, nil, err
	}

	err = page.SetContent(string(renderedTemplate))
	if err != nil {
		return nil, nil, err
	}

	data, err := page.Screenshot(playwright.PageScreenshotOptions{
		// Path: TODO(zikani03): Add option to save file to disk?
		FullPage:   playwright.Bool(true),
		Animations: playwright.ScreenshotAnimationsDisabled,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not create screenshot: %v", err)
	}

	defer browser.Close()
	defer pw.Stop()

	var options models.PrintOptions
	return data, &options, nil
}
