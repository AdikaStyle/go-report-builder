package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/nndi-oss/greypot/exporter"
	"github.com/palantir/stacktrace"
)

type reportService struct {
	pdf             exporter.ReportExporter
	png             exporter.ReportExporter
	templateService TemplateService
	baseUrl         string
}

func (rs *reportService) ExportReportHtml(reportId string, data interface{}) ([]byte, error) {
	html, err := rs.templateService.RenderTemplate(reportId, data)
	if err != nil {
		return nil, stacktrace.RootCause(err)
	}

	b64Html := base64.StdEncoding.EncodeToString(html)

	return []byte(b64Html), nil
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

	b64Pdf := base64.StdEncoding.EncodeToString(pdf)

	return []byte(b64Pdf), nil
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

	b64Png := base64.StdEncoding.EncodeToString(png)

	return []byte(b64Png), nil
}

func (rs *reportService) buildUrl(reportId string, data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", stacktrace.Propagate(err, "failed to marshal %+v to json", data)
	}

	b64Data := base64.StdEncoding.EncodeToString(jsonData)
	if b64Data == "e30=" {
		b64Data = ""
	}

	return fmt.Sprintf("http://%s/reports/render/%s?d=%s", rs.baseUrl, reportId, b64Data), nil
}

func NewReportService(pdf exporter.ReportExporter, png exporter.ReportExporter, templateService TemplateService, baseUrl string) *reportService {
	return &reportService{pdf: pdf, png: png, templateService: templateService, baseUrl: baseUrl}
}
