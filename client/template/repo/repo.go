package repo

type TemplateRepository interface {
	ListAll() ([]string, error)
	LoadTemplate(templateId string) ([]byte, error)
}
