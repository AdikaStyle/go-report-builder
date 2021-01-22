package business

type TemplateService interface {
	ListTemplates() ([]string, error)
	GetTemplate(reportId string) ([]byte, error)
	RenderTemplate(reportId string, data interface{}) ([]byte, error)
}

type ReportService interface {
	ExportReportHtml(reportId string, data interface{}) ([]byte, error)
	ExportReportPdf(reportId string, data interface{}) ([]byte, error)
	ExportReportPng(reportId string, data interface{}) ([]byte, error)
}
