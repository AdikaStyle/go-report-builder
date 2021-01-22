package data

import (
	"context"
	"github.com/AdikaStyle/go-report-builder/internal/models"
	"github.com/chromedp/chromedp"
	"github.com/palantir/stacktrace"
	"log"
	"time"
)

type pngReportExporter struct {
	timeout time.Duration
}

func NewPngReportExporter(timeout time.Duration) *pngReportExporter {
	return &pngReportExporter{timeout: timeout}
}

func (pre *pngReportExporter) Export(url string) ([]byte, *models.PrintOptions, error) {
	baseCtx, to := context.WithTimeout(context.Background(), pre.timeout)
	ctx, cancel := chromedp.NewContext(baseCtx)
	defer to()
	defer cancel()

	var res []byte
	var options models.PrintOptions
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Evaluate(extractPrintOptions, &options),
		chromedp.WaitVisible(`#printable`, chromedp.ByID),
		chromedp.Screenshot("#printable", &res, chromedp.NodeVisible, chromedp.ByID),
	)
	if err != nil {
		log.Println(err.Error())
		return nil, nil, stacktrace.RootCause(err)
	}

	return res, &options, nil
}

const extractPrintOptions = `function extractStyles() {
	var styles = getComputedStyle(document.body);

	return {
		"page_height": document.body.offsetHeight,
		"page_width": document.body.offsetWidth,
		"orientation": styles['orientation'] || "portrait",
	}
}
extractStyles();`
