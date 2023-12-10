package backend

import (
	"github.com/influenzanet/survey-repository/pkg/models"
)

type SurveyFilter struct {
	Platforms  []string
	ImporterAt RangeFilter
	Published  RangeFilter
	Limit      int
	Offset     int
}

type RangeFilter struct {
	From int64
	To   int64
}

type Backend interface {
	Start() error
	ImportSurvey(meta models.SurveyMetadata, data []byte) (uint, error)
	GetSurveys(namespace uint, filters SurveyFilter) ([]models.SurveyMetadata, error)
	GetNamespaces() ([]models.Namespace, error)
	CreateNamespace(name string) (uint, error)
	GetSurveyData(id uint, decompress bool) ([]byte, error)
	GetSurveyMeta(id uint) (models.SurveyMetadata, error)
	Migrate() error
}
