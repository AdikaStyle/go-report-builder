package business

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/AdikaStyle/go-report-builder/internal/data"
	"github.com/palantir/stacktrace"
)

type reportService struct {
	pdf             data.ReportExporter
	png             data.ReportExporter
	templateService TemplateService
	baseUrl         string
}

func (rs *reportService) ExportReportHtml(reportId string, data interface{}) ([]byte, error) {
	html, err := rs.templateService.RenderTemplate(reportId, data)
	if err != nil {
		return nil, stacktrace.RootCause(err)
	}

	var b64Html []byte
	base64.StdEncoding.Encode(b64Html, html)

	return b64Html, nil
}

func (rs *reportService) ExportReportPdf(reportId string, data interface{}) ([]byte, error) {
	url, err := rs.buildUrl(reportId, data)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to build url for pdf export on reportId: %s and data: %+v", reportId, data)
	}

	pdf, _, err := rs.pdf.Export(url)
	if err != nil {
		return nil, stacktrace.RootCause(err)
	}

	var b64Pdf []byte
	base64.StdEncoding.Encode(b64Pdf, pdf)

	return b64Pdf, nil
}

func (rs *reportService) ExportReportPng(reportId string, data interface{}) ([]byte, error) {
	url, err := rs.buildUrl(reportId, data)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to build url for pdf export on reportId: %s and data: %+v", reportId, data)
	}

	png, _, err := rs.png.Export(url)
	if err != nil {
		return nil, stacktrace.RootCause(err)
	}

	var b64Png []byte
	base64.StdEncoding.Encode(b64Png, png)

	return b64Png, nil
}

func (rs *reportService) buildUrl(reportId string, data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", stacktrace.Propagate(err, "failed to marshal %+v to json", data)
	}

	var b64Data []byte
	base64.StdEncoding.Encode(b64Data, jsonData)

	return fmt.Sprintf("http://%s/reports/%s?d=%s", rs.baseUrl, reportId, b64Data), nil
}

func NewReportService(pdf data.ReportExporter, png data.ReportExporter, templateService TemplateService, baseUrl string) *reportService {
	return &reportService{pdf: pdf, png: png, templateService: templateService, baseUrl: baseUrl}
}
