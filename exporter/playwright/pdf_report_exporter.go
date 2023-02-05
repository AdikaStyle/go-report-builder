package playwright

import (
	"fmt"
	"io/ioutil"
	"os"
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

	tmpFile, err := os.CreateTemp(os.TempDir(), ".pdf")

	if err != nil {
		return nil, nil, err
	}

	_, err = page.PDF(playwright.PagePdfOptions{
		Path: playwright.String(tmpFile.Name()),
	})

	if err != nil {
		return nil, nil, err
	}

	data, err := ioutil.ReadAll(tmpFile)

	if err != nil {
		return nil, nil, err
	}

	defer tmpFile.Close()
	defer browser.Close()
	defer pw.Stop()

	var options models.PrintOptions
	return data, &options, nil
}
