package surveys

// SurveyV1_3 structure for Survey before V1_3
type SurveyV1_2 struct {
	ID      string            `json:"id"`
	Current SurveyVersionV1_2 `json:"current"`
}

type SurveyVersionV1_2 struct {
	Published        int64          `json:"published,string"`
	UnPublished      int64          `json:"unpublished,string"`
	SurveyDefinition SurveyItemV1_2 `json:"surveyDefinition"`
	VersionID        string         `json:"versionID"`
}

type SurveyItemV1_2 struct {
	Key string `json:"key"`
}

type EmbedddedSurveyV1_2 struct {
	StudyKey string     `json:"studyKey"`
	Survey   SurveyV1_2 `json:"survey"`
}
