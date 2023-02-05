package chromedp

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/nndi-oss/greypot/models"
	"github.com/palantir/stacktrace"
)

type pngReportExporter struct {
	timeout        time.Duration
	viewportHeight int
	viewportWidth  int
}

func NewPngReportExporter(timeout time.Duration, vpHeight, vpWidth int) *pngReportExporter {
	return &pngReportExporter{timeout: timeout, viewportHeight: vpHeight, viewportWidth: vpWidth}
}

func (pre *pngReportExporter) Export(url string, renderedTemplate []byte) ([]byte, *models.PrintOptions, error) {
	ctx, cancel := createContext(pre.timeout, pre.viewportHeight, pre.viewportWidth)
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

func createContext(timeout time.Duration, vph int, vpw int) (context.Context, context.CancelFunc) {
	baseCtx, cancelTimeout := context.WithTimeout(context.Background(), timeout)

	opts := []chromedp.ExecAllocatorOption{
		chromedp.WindowSize(vpw, vph),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
	}
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(baseCtx, opts...)
	ctx, cancelCtx := chromedp.NewContext(allocCtx)

	cancelFuncs := func() {
		cancelTimeout()
		cancelAlloc()
		cancelCtx()
	}

	return ctx, cancelFuncs
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
