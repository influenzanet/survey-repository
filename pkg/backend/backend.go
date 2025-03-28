package backend

import (
	"github.com/influenzanet/survey-repository/pkg/models"
)

type SurveyFilter struct {
	Platforms  []string
	Names      []string // Survey names
	ModelTypes []string // Model types
	ImporterAt RangeFilter
	Published  RangeFilter
	Limit      int
	Offset     int
}

type RangeFilter struct {
	From int64
	To   int64
}

type PaginateInfo struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
	Total  int64 `json:"total"`
}

type PaginatedResult[T any] struct {
	PaginateInfo
	Data []T `json:"data"`
}

type Backend interface {
	Start() error
	FindSurvey(meta models.SurveyMetadata) (uint, error)
	ImportSurvey(meta models.SurveyMetadata, data []byte) (uint, error)
	GetSurveys(namespace uint, filters SurveyFilter) (PaginatedResult[models.SurveyMetadata], error)
	GetNamespaces() ([]models.Namespace, error)
	CreateNamespace(name string) (uint, error)
	GetSurveyData(id uint, decompress bool) ([]byte, error)
	GetSurveyMeta(id uint) (models.SurveyMetadata, error)
	Migrate() error
}
