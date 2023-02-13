package exporter

import (
	"time"

	cd "github.com/nndi-oss/greypot/exporter/chromedp"
	pw "github.com/nndi-oss/greypot/exporter/playwright"
	"github.com/nndi-oss/greypot/models"
)

type ReportExporter interface {
	Export(url string, renderedTemplate []byte) ([]byte, *models.PrintOptions, error)
}

func NewChromePNGReportExporter(timeout time.Duration, vpHeight, vpWidth int) ReportExporter {
	return cd.NewPngReportExporter(timeout, vpHeight, vpWidth)
}

func NewChromePDFReportExporter(timeout time.Duration, vpHeight, vpWidth int) ReportExporter {
	return cd.NewPdfReportExporter(timeout, vpHeight, vpWidth)
}

func NewPlaywrightPDFReportExporter(timeout time.Duration, vpHeight, vpWidth int) ReportExporter {
	return pw.NewPdfReportExporter(timeout)
}

func NewPlaywrightPNGReportExporter(timeout time.Duration, vpHeight, vpWidth int) ReportExporter {
	return pw.NewPngReportExporter(timeout)
}
