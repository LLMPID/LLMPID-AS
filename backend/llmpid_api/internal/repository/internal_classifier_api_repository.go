package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"llm-promp-inj.api/internal/dto"
)

type InternalClassifierAPIRepository struct {
	apiPath string
	logger  *logrus.Logger
}

func NewInternalClassifierAPIRepository(apiPath string, logger *logrus.Logger) *InternalClassifierAPIRepository {
	return &InternalClassifierAPIRepository{apiPath: apiPath, logger: logger}
}

func (r *InternalClassifierAPIRepository) SendClassificationRequest(classificationRequest dto.ClassificationRequest) (string, error) {
	// Marshal the classification request to JSON.
	requestBody, err := json.Marshal(classificationRequest)
	if err != nil {
		r.logger.Error("Uanble to serialize the classification request into JSON before sending it to the internal classification service API. ERR: ", err)
		return "", err
	}

	// Make the classification POST request to the internal classification service.
	resp, err := http.Post(r.apiPath, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		r.logger.Error("Unable to perform request to the internal classification service API. ERR: ", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check for non-200 HTTP response codes in case the classification failed.
	if resp.StatusCode != http.StatusOK {
		r.logger.Error("The internal classification service was unable to classify the request. Response status code: ", resp.StatusCode)
		return "", errors.New("failed to send request: " + resp.Status)
	}

	// Read and parse the response to validate its integrity.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.logger.Error("Unable to process raw response body from the internal classification service. ERR: ", resp.StatusCode)
		return "", err
	}

	// Parse the response body into a map in order to access the "result" field.
	var resultMap map[string]interface{}
	if err := json.Unmarshal(body, &resultMap); err != nil {
		r.logger.Error("Reesponse body from the internal classification service in not a valid JSON. ERR: ", resp.StatusCode)
		return "", err
	}

	// Access the "result" field directly and get the classification result. Can be either "Injection" or "Normal".
	result, ok := resultMap["result"].(string)
	if !ok {
		r.logger.Error("Reesponse body from the internal classification service does not contain result field. ERR: ", resp.StatusCode)
		return "", errors.New("result of internal classification API request is not a string")
	}

	return result, nil
}
