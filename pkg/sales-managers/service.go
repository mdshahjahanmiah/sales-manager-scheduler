package sales_managers

import (
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/sales-manager-scheduler/pkg/config"
	"github.com/mdshahjahanmiah/sales-manager-scheduler/pkg/db"
)

type Service interface {
	AvailableSlots(request QueryRequest) ([]AvailableSlot, error)
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

func (s service) AvailableSlots(request QueryRequest) ([]AvailableSlot, error) {
	// Query matching sales managers
	managers, err := FindMatchingSalesManagers(s.db.DB, request)
	if err != nil {
		return []AvailableSlot{}, err
	}

	// If no matching sales managers found, return an empty response
	if len(managers) == 0 {
		return []AvailableSlot{}, nil
	}

	s.logger.Info("managers information", "managers", managers)

	// Query for available slots
	availableSlots, err := FindAvailableSlots(s.db.DB, request, managers)
	if err != nil {
		return []AvailableSlot{}, err
	}

	return availableSlots, nil

}
