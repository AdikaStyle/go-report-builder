package exporter

import "github.com/nndi-oss/greypot/models"

type ReportExporter interface {
	Export(url string) ([]byte, *models.PrintOptions, error)
}
