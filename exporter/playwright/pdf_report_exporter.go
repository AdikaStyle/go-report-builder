package playwright

import (
	"fmt"
	"time"

	"github.com/nndi-oss/greypot/models"
	"github.com/playwright-community/playwright-go"
)

type pdfReportExporter struct {
	timeout time.Duration
	// viewportHeight int
	// viewportWidth  int
}

func NewPdfReportExporter(timeout time.Duration) *pdfReportExporter {
	return &pdfReportExporter{timeout: timeout}
}

func (pre *pdfReportExporter) Export(url string, renderedTemplate []byte) ([]byte, *models.PrintOptions, error) {
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

	data, err := page.PDF(playwright.PagePdfOptions{
		// Path: TODO(zikani03): Add option to save file to disk?
	})
	if err != nil {
		return nil, nil, err
	}

	defer browser.Close()
	defer pw.Stop()

	var options models.PrintOptions
	return data, &options, nil
}
