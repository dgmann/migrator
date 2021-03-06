package databasereader

import (
	"path"
	"github.com/dgmann/document-manager/migrator/records/models"
)

type PdfFile struct {
	models.SubRecord
	PathExtension *string    `db:"PathExtension"`
}

func (p *PdfFile) AsSubRecord() *models.SubRecord {
	s := p.SubRecord
	s.Path = p.GetPath()
	return &s
}

func (p *PdfFile) GetPath() string {
	if p.PathExtension == nil {
		return p.Name
	}
	return path.Join(*p.PathExtension, p.Name)
}
