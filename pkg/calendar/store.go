package calendar

import (
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/mdshahjahanmiah/sales-manager-scheduler/pkg/db"
	"strings"
	"time"
)

// Store defines the interface for matching sales managers and finding available slots.
type Store interface {
	GetSalesManagers(req queryRequest) ([]int, error)
	GetAvailableSlots(req queryRequest, managers []int) ([]AvailableSlot, error)
}

// slot is a struct used for storing slot information
type slot struct {
	SalesManagerID int
	StartDate      time.Time
}

// store holds the database connection and implements the Store interface.
type store struct {
	db *db.DB
}

// NewStore creates a new store instance with the provided database connection.
func NewStore(db *db.DB) *store {
	return &store{db: db}
}

// GetSalesManagers retrieves the IDs of sales managers based on the specified query parameters.
func (s *store) GetSalesManagers(req queryRequest) ([]int, error) {
	query := `
        SELECT id
        FROM sales_managers
        WHERE $1 = ANY(languages) -- Check if the manager's languages include the requested language.
        AND products @> $2  -- Check if the manager's products include all the requested products.
        AND $3 = ANY(customer_ratings) -- Check if the manager's customer ratings include the requested rating.
    `
	rows, err := s.db.DB.Query(query, req.Language, pq.Array(req.Products), req.Rating)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var managerIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		managerIDs = append(managerIDs, id)
	}

	return managerIDs, nil
}

// GetAvailableSlots retrieves the available slots for the given sales managers on a specified date.
// It combines information about booked and free slots, checks for conflicts, and returns the available slots
func (s *store) GetAvailableSlots(req queryRequest, managers []int) ([]AvailableSlot, error) {
	managerIDs := strings.Trim(strings.Replace(fmt.Sprint(managers), " ", ",", -1), "[]")

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	// Fetch booked slots for the given managers on the specified date.
	bookedSlots, err := s.fetchBookedSlots(managerIDs, date)
	if err != nil {
		return nil, err
	}

	// Fetch free slots for the given managers on the specified date.
	freeSlots, err := s.fetchFreeSlots(managerIDs, date)
	if err != nil {
		return nil, err
	}

	// Calculate available slots by checking the free slots against the booked slots for conflicts.
	slotAvailability := s.calculateAvailableSlots(freeSlots, bookedSlots)

	return s.mapToAvailableSlots(slotAvailability), nil
}

// fetchBookedSlots retrieves booked slots for the specified sales managers on a given date.
func (s *store) fetchBookedSlots(managerIDs string, date time.Time) (map[int][][2]time.Time, error) {
	query := `
    SELECT sales_manager_id, start_date, end_date, booked
    FROM slots
    WHERE sales_manager_id = ANY($1) -- Filters slots belonging to the provided manager IDs.
    AND start_date::date = $2 -- Filters slots that match the specified date.
`
	rows, err := s.db.DB.Query(query, "{"+managerIDs+"}", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookedSlots := make(map[int][][2]time.Time)
	for rows.Next() {
		var salesManagerID int
		var startDate, endDate time.Time
		var booked bool
		if err := rows.Scan(&salesManagerID, &startDate, &endDate, &booked); err != nil {
			return nil, err
		}

		// Only add the slot to the bookedSlots map if it is marked as booked.
		if booked {
			// Append the booked time range (start and end times) to the slice for the corresponding sales manager.
			bookedSlots[salesManagerID] = append(bookedSlots[salesManagerID], [2]time.Time{startDate, endDate})
		}
	}

	return bookedSlots, nil
}

// fetchFreeSlots retrieves the free (unbooked) slots for the specified sales managers on a given date.
func (s *store) fetchFreeSlots(managerIDs string, date time.Time) ([]slot, error) {
	query := `
    SELECT sales_manager_id, start_date
    FROM slots
    WHERE sales_manager_id = ANY($1)  -- Filters slots belonging to the provided manager IDs.
    AND start_date::date = $2  -- Filters slots that match the specified date.
    AND booked = FALSE  -- Only selects slots that are not booked.
`
	rows, err := s.db.DB.Query(query, "{"+managerIDs+"}", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var freeSlots []slot
	for rows.Next() {
		var slot slot
		if err := rows.Scan(&slot.SalesManagerID, &slot.StartDate); err != nil {
			return nil, err
		}
		// Append the retrieved slot to the freeSlots slice.
		freeSlots = append(freeSlots, slot)
	}

	return freeSlots, nil
}

// calculateAvailableSlots determines the availability of free slots by checking them against booked slots for conflicts.
func (s *store) calculateAvailableSlots(freeSlots []slot, bookedSlots map[int][][2]time.Time) map[string]int {
	slotAvailability := make(map[string]int)

	for _, freeSlot := range freeSlots {
		// Check if the current free slot conflicts with any of the booked slots for the same sales manager.
		conflict := s.hasConflict(freeSlot, bookedSlots[freeSlot.SalesManagerID])
		if !conflict {
			// If there's no conflict, format the start date of the slot to a string and increment its availability count.
			formattedDate := freeSlot.StartDate.Format("2006-01-02T15:04:05.000Z")
			slotAvailability[formattedDate]++
		}
	}

	return slotAvailability
}

// hasConflict checks whether a free slot conflicts with any of the booked slots for a given sales manager.
func (s *store) hasConflict(freeSlot slot, bookedSlots [][2]time.Time) bool {
	for _, bSlot := range bookedSlots {
		// Check if the start time of the free slot is before the end time of a booked slot
		// and the end time of the free slot (one hour after start) is after the start time of the booked slot.
		// If both conditions are met, the slots overlap, indicating a conflict.
		if freeSlot.StartDate.Before(bSlot[1]) && freeSlot.StartDate.Add(time.Hour).After(bSlot[0]) {
			// Return true if a conflict is detected.
			return true
		}
	}
	// Return false if no conflicts are found with any booked slots.
	return false
}

// mapToAvailableSlots converts the slot availability map into a slice of AvailableSlot structs.
func (s *store) mapToAvailableSlots(slotAvailability map[string]int) []AvailableSlot {
	var slots []AvailableSlot
	for startDate, count := range slotAvailability {
		// Append a new AvailableSlot struct to the slots slice with the start date and availability count.
		slots = append(slots, AvailableSlot{
			AvailableCount: count,     // Number of available slots for this start date.
			StartDate:      startDate, // Start date of the available slot formatted as a string.
		})
	}
	return slots
}
