package surveys

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/influenzanet/survey-repository/pkg/models"
)

var ErrUnknownSurveyModel = errors.New("unknown survey model")
var ErrWrongSurveyModel = errors.New("wrong survey model")
var ErrUnexpectedEntryType = errors.New("wrong format, unexpected type")

func computeSha(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// Try to find out the survey version and some metadata
func ExtractSurveyMetadata(data []byte) (*models.SurveyDescriptor, error) {
	d, err := detectSurveyMetadata(data)
	if err != nil {
		return nil, err
	}
	d.Sha = computeSha(data)
	return d, nil
}

func detectSurveyMetadata(data []byte) (*models.SurveyDescriptor, error) {
	var err error
	var d *models.SurveyDescriptor

	d, err = tryVersion1_3(data)

	if err == nil {
		fmt.Println("Survey v1.3")
		fmt.Println(d)
		return d, nil
	}

	if !errors.Is(err, ErrWrongSurveyModel) {
		return nil, err
	}

	d, err = tryVersion1_2(data)
	if err == nil {
		return d, nil
	}

	d, err = tryVersion1_2_embeded(data)
	if err == nil {
		return d, nil
	}

	return nil, ErrUnknownSurveyModel
}

func tryVersion1_2_embeded(data []byte) (*models.SurveyDescriptor, error) {
	s := EmbedddedSurveyV1_2{}
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}
	return fromVersion1_2(s.Survey)
}

func tryVersion1_2(data []byte) (*models.SurveyDescriptor, error) {
	s := SurveyV1_2{}
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}
	return fromVersion1_2(s)
}

func fromVersion1_2(s SurveyV1_2) (*models.SurveyDescriptor, error) {
	if s.Current.SurveyDefinition.Key == "" {
		return nil, errors.Join(ErrWrongSurveyModel, errors.New("SurveyDefinition.key is empty"))
	}
	d := models.SurveyDescriptor{}
	d.Name = s.Current.SurveyDefinition.Key
	d.VersionID = s.Current.VersionID
	d.ExternalID = s.ID
	d.Published = s.Current.Published
	d.ModelVersion = models.SurveyVersion1_2
	return &d, nil
}

func tryVersion1_3(data []byte) (*models.SurveyDescriptor, error) {
	s := SurveyV1_3{}
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}
	if s.SurveyDefinition.Key == "" {
		return nil, errors.Join(ErrWrongSurveyModel, errors.New("SurveyDefinition.key is empty"))
	}

	d := models.SurveyDescriptor{}
	d.Name = s.SurveyDefinition.Key
	d.VersionID = s.VersionID
	d.ExternalID = s.ID
	d.Published = s.Published
	d.ModelVersion = models.SurveyVersion1_3
	return &d, nil
}
