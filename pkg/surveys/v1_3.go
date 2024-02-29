package surveys

type SurveyV1_3 struct {
	ID               string               `json:"id,omitempty"`
	Published        int64                `json:"published,string"`
	Unpublished      int64                `json:"unpublished,string"`
	SurveyDefinition SurveyDefinitionV1_3 `json:"surveyDefinition"`
	VersionID        string               `json:"versionID"`
	Metadata         map[string]string    `json:"metadata,omitempty"`
}

type SurveyDefinitionV1_3 struct {
	Key string `json:"key"`
}
