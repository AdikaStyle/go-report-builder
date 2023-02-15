package exporter

import (
	"time"

	pw "github.com/nndi-oss/greypot/exporter/playwright"
	"github.com/nndi-oss/greypot/models"
)

type ReportExporter interface {
	Export(url string, renderedTemplate []byte) ([]byte, *models.PrintOptions, error)
}

func NewPlaywrightPDFReportExporter(timeout time.Duration, vpHeight, vpWidth int) ReportExporter {
	return pw.NewPdfReportExporter(timeout)
}

func NewPlaywrightPNGReportExporter(timeout time.Duration, vpHeight, vpWidth int) ReportExporter {
	return pw.NewPngReportExporter(timeout)
}
