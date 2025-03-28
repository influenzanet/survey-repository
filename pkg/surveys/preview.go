package surveys

type SurveyPreview struct {
	VersionID string `json:"versionId"`
	Published int64 `json:"published"`
	Questions map[string]interface{} `json:"questions"`
}

type SurveyPreviewBundle struct {
	Key  string `json:"key"`
	Versions []SurveyPreview `json:"versions"`
}
