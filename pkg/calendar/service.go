package calendar

import (
	"errors"
	eError "github.com/mdshahjahanmiah/explore-go/error"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/sales-manager-scheduler/pkg/config"
	"github.com/mdshahjahanmiah/sales-manager-scheduler/pkg/db"
	"net/http"
)

type Service interface {
	AvailableSlots(request queryRequest) ([]AvailableSlot, error)
}

type AvailableSlot struct {
	AvailableCount int    `json:"available_count"`
	StartDate      string `json:"start_date"`
}

type service struct {
	config config.Config
	logger *logging.Logger
	store  *store
}

// NewService creates a new calendar service instance
func NewService(config config.Config, logger *logging.Logger, database *db.DB) (Service, error) {
	return &service{
		config: config,
		logger: logger,
		store:  NewStore(database),
	}, nil
}

// AvailableSlots retrieves available slots for sales managers that match the given request criteria.
func (s service) AvailableSlots(request queryRequest) ([]AvailableSlot, error) {
	// Get matching sales managers based on the request criteria (e.g., language, products, ratings).
	managers, err := s.store.GetSalesManagers(request)
	if err != nil {
		s.logger.Error("error matching sales managers", "err", err.Error())
		return []AvailableSlot{}, eError.NewTransportError(errors.New("please try again later"), "internal_server_error")
	}

	// If no matching sales managers are found, return an empty response with a validation error.
	if len(managers) == 0 {
		return []AvailableSlot{}, eError.NewServiceError(errors.New("no matching sales managers found"), "validation_error", "payload", http.StatusBadRequest)
	}

	// Get available slots for the matching sales managers based on the request criteria.
	availableSlots, err := s.store.GetAvailableSlots(request, managers)
	if err != nil {
		s.logger.Error("error getting available slots by sales managers", "err", err.Error())
		return []AvailableSlot{}, eError.NewTransportError(errors.New("please try again later"), "internal_server_error")
	}

	// If no available slots are found, return an empty slice with no error.
	if len(availableSlots) == 0 {
		return []AvailableSlot{}, nil
	}

	return availableSlots, nil

}
