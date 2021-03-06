package filesystem

import (
	"strings"
	"path/filepath"
	"github.com/pkg/errors"
	"strconv"
	"os"
	"github.com/dgmann/document-manager/migrator/records/models"
	pdf "github.com/unidoc/unidoc/pdf/model"
	"github.com/sirupsen/logrus"
)

type EmbeddedRecord = models.Record

type Record struct {
	*EmbeddedRecord
	SplittedPdfDir string
}

func NewRecordFromPath(dir string) (*Record, error) {
	directory, fileName := filepath.Split(dir)
	if directory == "" {
		return nil, errors.New("file dir cannot be parsed to create a record: " + dir)
	}

	spezialization := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	patIdString := filepath.Base(directory)
	patId, err := strconv.Atoi(patIdString)
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert patId to integer: "+dir)
	}

	r := &models.Record{
		Name:  fileName,
		Path:  dir,
		Spez:  spezialization,
		PatId: patId,
		Pages: -2,
	}
	return &Record{EmbeddedRecord: r}, nil
}

func (r *Record) PageCount() int {
	if r.Pages != -2 {
		return r.Pages
	}
	pageCount, err := getPageCount(r.Path)
	if err != nil {
		logrus.Warn("cannot read page count of file: ", r.Path)
	}
	r.Pages = pageCount
	return r.Pages
}

func getPageCount(file string) (int, error) {
	f, err := os.Open(file)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return -1, err
	}

	return pdfReader.GetNumPages()
}

func (r *Record) GetPath() string {
	return r.Path
}

func (r *Record) Close() error {
	return os.RemoveAll(r.SplittedPdfDir)
}
