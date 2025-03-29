package models

type SurveyStats struct {
	PlatformID string `json:"platform"`
	ModelType  string `json:"model_type"`
    DescriptorName	   string `json:"survey_key" gorm:"descriptor_name"`
	Count      int    `json:"count"`
}