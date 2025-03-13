package service

import (
	"fmt"

	"llm-promp-inj.api/internal/dto"
	"llm-promp-inj.api/internal/models"
	"llm-promp-inj.api/internal/repository"
)

type ClassificationRequestService struct {
	ClassificationLogsRepo *repository.ClassificationLogRepository
	ClassificationRepo     *repository.InternalClassifierAPIRepository
}

func NewClassificationRequestService(logsRepo *repository.ClassificationLogRepository, clsRepo *repository.InternalClassifierAPIRepository) *ClassificationRequestService {
	return &ClassificationRequestService{
		ClassificationLogsRepo: logsRepo,
		ClassificationRepo:     clsRepo,
	}
}

// ClassifyText performs prompt injection classification for a privded string.
// It first sends the string for classification to an internal API, retrieves and logs the result into a database, and then returns it to the client service.
func (s *ClassificationRequestService) ClassifyText(ClassificationLog dto.ClassificationRequest) (*models.ClassificationLog, error) {
	// Send data for classification.
	clssResult, err := s.ClassificationRepo.SendClassificationRequest(ClassificationLog)
	if err != nil {
		return nil, err
	}

	// Create a classification log with the request and result and make a DB entry.
	clssRequest := &models.ClassificationLog{RequestText: ClassificationLog.Text, Result: clssResult}
	err = s.ClassificationLogsRepo.InsertClassificationLog(clssRequest)
	if err != nil {
		return nil, err
	}

	return clssRequest, nil
}

func (s *ClassificationRequestService) GetClassificationRequestLogByID(id uint) (*models.ClassificationLog, error) {
	clssRequest, err := s.ClassificationLogsRepo.SelectClassificationLogByID(id)
	if err != nil {
		return nil, err
	}
	return clssRequest, nil
}

func (s *ClassificationRequestService) GetClassificationLogsByPage(page int, limit int, sortBy string) (*[]models.ClassificationLog, error) {
	var orderBy string

	// Assures that the sortBy parameter is valid. Defaults to "desc" if it is not.
	switch sortBy {
	case "desc", "asc":
		break
	default:
		sortBy = "desc"
	}

	// Creates orderBy parameter for the database query.
	orderBy = fmt.Sprintf("id %s", sortBy)

	clssRequests, err := s.ClassificationLogsRepo.SelectClassificationLogsByPage(page, limit, orderBy)
	if err != nil {
		return nil, err
	}

	return clssRequests, nil
}
