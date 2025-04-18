package models

type SurveyDescriptor struct {
	Name         string `json:"name" gorm:"index"`
	VersionID    string `json:"version"`
	ExternalID   string `json:"external_id"`
	Published    int64  `json:"published"`
	ModelVersion string `json:"model_version"` // Survey Model version
	Sha          string `json:"sha256"`        // Base64Url encoded sha256
}

const (
	SurveyVersion1_2 = "1.2" // Survey Model for study service for version 1.2 (before 1.3)
	SurveyVersion1_3 = "1.3" // Survey Model for study service for version 1.3+
	SurveyVersionPreview = "preview"

	SurveyModelDefinition = "D"
	SurveyModelPreview = "P"
)

type DBId struct {
	ID uint
}

type SurveyMetadata struct {
	ID         uint              `json:"id" gorm:"primaryKey"`
	Namespace  uint              `json:"namespace" gorm:"index"`
	ImportedAt int64             `json:"imported_at"`
	ImportedBy string            `json:"imported_by"`
	PlatformID string            `json:"platform" gorm:"index"`
	Version    string 			 `json:"version"`
	ModelType  string 			 `json:"model_type"` // Model type 'definition','preview'
	Labels     map[string]string `json:"labels" gorm:"serializer:json"`
	Descriptor SurveyDescriptor  `json:"descriptor" gorm:"embedded;embeddedPrefix:descriptor_"`
	SurveyData SurveyData        `json:"-" gorm:"foreignKey:SurveyID"` // Do not serialize this field
}

type SurveyData struct {
	SurveyID uint
	Survey   string
}

type Namespace struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"uniqueIndex"`
}
