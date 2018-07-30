package filesystem

import (
	"io"
	"github.com/dgmann/document-manager/migrator/records/models"
	"path/filepath"
	"github.com/dgmann/document-manager/migrator/splitter"
	"github.com/dgmann/document-manager/migrator/shared"
	"encoding/gob"
)

type Index struct {
	models.Index
	Path string
}

func newIndex(data []RecordContainerCloser, path string) *Index {
	var cast []models.RecordContainer
	for _, d := range data {
		cast = append(cast, d)
	}
	return &Index{Index: *models.NewIndex("filesystem", cast), Path: filepath.Clean(path)}
}

func (i *Index) Destroy() {
	for _, r := range i.Records() {
		if closer, ok := r.(io.Closer); ok {
			closer.Close()
		}
	}
}

func (i *Index) Validate() []string {
	invalidDirectories := make(map[string]struct{})
	for _, r := range i.Records() {
		dir := filepath.Dir(r.Record().Path)
		d := filepath.Dir(dir)
		if d != i.Path {
			invalidDirectories[dir] = struct{}{}
		}
	}
	keys := make([]string, 0, len(invalidDirectories))
	for k := range invalidDirectories {
		keys = append(keys, k)
	}
	return keys
}

func (i *Index) LoadSubRecords(dir string) error {
	var err error
	for _, r := range i.Records() {
		e := loadSubRecord(r.(*Record), dir)
		if e != nil {
			err = shared.WrapError(err, e.Error())
		}
	}
	return err
}

func loadSubRecord(record *Record, dir string) error {
	if len(record.SubRecords) > 0 { // Already loaded
		return nil
	}
	subrecords, tmpDir, err := splitter.Split(record.Path, dir)
	if err != nil {
		return err
	}
	var convertedSubRecords []models.SubRecordContainer
	for _, subrecord := range subrecords {
		subrecord.BefundId = &record.Id
		subrecord.PatId = &record.PatId
		subrecord.Spezialization = &record.Spez
		convertedSubRecords = append(convertedSubRecords, &SubRecord{*subrecord})
	}
	record.SubRecords = convertedSubRecords
	record.SplittedPdfDir = tmpDir
	return nil
}

func (i *Index) SubRecords() []models.SubRecordContainer {
	var subrecords []models.SubRecordContainer
	for _, record := range i.Records() {
		subrecords = append(subrecords, record.Record().SubRecords...)
	}
	return subrecords
}

func (i *Index) Save(dir string) error {
	gob.Register(&Record{})
	gob.Register(&SubRecord{})
	return i.Index.Save(dir)
}

func (i *Index) Load(path string) error {
	gob.Register(&Record{})
	gob.Register(&SubRecord{})
	return i.Index.Load(path)
}
