package repo

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type fsTemplateRepo struct {
	root            string
	templatesFolder fs.FS
}

func NewFSTemplateRepo(templatesFolder fs.FS) *fsTemplateRepo {
	return &fsTemplateRepo{
		root:            ".",
		templatesFolder: templatesFolder,
	}
}

func (ftr *fsTemplateRepo) ListAll() ([]string, error) {
	return ftr.listFiles(ftr.templatesFolder)
}

func (ftr *fsTemplateRepo) LoadTemplate(templateId string) ([]byte, error) {
	return fs.ReadFile(ftr.templatesFolder, filepath.Join(ftr.root, templateId+".html"))
}

func (ftr *fsTemplateRepo) listFiles(root fs.FS) ([]string, error) {
	var files []string

	err := fs.WalkDir(ftr.templatesFolder, ftr.root, func(path string, info fs.DirEntry, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			key := strings.ReplaceAll(strings.TrimPrefix(strings.TrimPrefix(path, ftr.root), "/"), ".html", "")
			files = append(files, key)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, err
}
