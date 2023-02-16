package main

import (
	"fmt"
	"strings"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type InmemTemplate struct {
	Name    string
	Content []byte
}

type inmemoryTemplateRepository struct {
	files cmap.ConcurrentMap[string, InmemTemplate]
}

func NewInmemoryTemplateRepository() *inmemoryTemplateRepository {
	return &inmemoryTemplateRepository{
		files: cmap.New[InmemTemplate](),
	}
}

func (ftr *inmemoryTemplateRepository) Add(name, content string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	data := InmemTemplate{
		Name:    strings.TrimSpace(name),
		Content: []byte(content),
	}

	ftr.files.Set(name, data)
	return nil
}

func (ftr *inmemoryTemplateRepository) ListAll() ([]string, error) {
	return ftr.listFiles("")
}

func (ftr *inmemoryTemplateRepository) LoadTemplate(templateId string) ([]byte, error) {
	item, found := ftr.files.Get(strings.TrimSpace(templateId))
	if !found {
		return nil, fmt.Errorf("template not found in in-mem ory repository")
	}
	return item.Content, nil
}

func (ftr *inmemoryTemplateRepository) listFiles(root string) ([]string, error) {
	var files []string

	for item := range ftr.files.IterBuffered() {
		val := item.Val
		if val.Name == "" {
			continue
		}
		key := strings.ReplaceAll(strings.TrimSpace(val.Name), ".html", "")
		files = append(files, key)
	}

	return files, nil
}
