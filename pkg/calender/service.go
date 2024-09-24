package calender

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
	db     *db.DB
}

func NewService(config config.Config, logger *logging.Logger, database *db.DB) (Service, error) {
	return &service{
		config: config,
		logger: logger,
		db:     database,
	}, nil
}

func (s service) AvailableSlots(request queryRequest) ([]AvailableSlot, error) {
	// get matching sales managers
	managers, err := matchingSalesManagers(s.db.DB, request)
	if err != nil {
		s.logger.Error("error matching sales managers", "err", err.Error())
		return []AvailableSlot{}, eError.NewTransportError(errors.New("please try again later"), "internal_server_error")
	}

	// If no matching sales managers found, return an empty response
	if len(managers) == 0 {
		return []AvailableSlot{}, eError.NewServiceError(errors.New("no matching sales managers found"), "validation_error", "payload", http.StatusBadRequest)
	}

	// get available slots by matching sales managers
	availableSlots, err := availableSlots(s.db.DB, request, managers)
	if err != nil {
		s.logger.Error("error getting available slots by sales managers", "err", err.Error())
		return []AvailableSlot{}, eError.NewTransportError(errors.New("please try again later"), "internal_server_error")
	}

	// If available slots found, return an empty response []
	if len(availableSlots) == 0 {
		return []AvailableSlot{}, nil
	}

	return availableSlots, nil

}
