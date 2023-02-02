package repo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type filesystemTemplateRepo struct {
	templatesFolder string
}

func NewFilesystemTemplateRepo(templatesFolder string) *filesystemTemplateRepo {
	return &filesystemTemplateRepo{templatesFolder: templatesFolder}
}

func (ftr *filesystemTemplateRepo) ListAll() ([]string, error) {
	return ftr.listFiles(ftr.templatesFolder)
}

func (ftr *filesystemTemplateRepo) LoadTemplate(templateId string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(ftr.templatesFolder, templateId+".html"))
}

func (ftr *filesystemTemplateRepo) listFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			key := strings.ReplaceAll(strings.TrimPrefix(strings.TrimPrefix(path, ftr.templatesFolder), "/"), ".html", "")
			files = append(files, key)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, err
}
