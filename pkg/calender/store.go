package calender

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strings"
	"time"
)

func matchingSalesManagers(db *sql.DB, req queryRequest) ([]int, error) {
	query := `
        SELECT id
        FROM sales_managers
        WHERE $1 = ANY(languages)
        AND products @> $2
        AND $3 = ANY(customer_ratings)
    `
	rows, err := db.Query(query, req.Language, pq.Array(req.Products), req.Rating)
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

func availableSlots(db *sql.DB, req queryRequest, managers []int) ([]AvailableSlot, error) {
	managerIDs := strings.Trim(strings.Replace(fmt.Sprint(managers), " ", ",", -1), "[]")

	// get all slots for the given sales managers and date
	query := `
      SELECT sales_manager_id, start_date, end_date, booked
      FROM slots
      WHERE sales_manager_id = ANY($1)
      AND start_date::date = $2
  `
	rows, err := db.Query(query, "{"+managerIDs+"}", req.Date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to hold the availability count for each slot
	slotAvailability := make(map[string]int)

	// Map to track booked slots for each sales manager
	bookedSlots := make(map[int][][2]time.Time)

	for rows.Next() {
		var salesManagerID int
		var startDate, endDate time.Time
		var booked bool
		if err := rows.Scan(&salesManagerID, &startDate, &endDate, &booked); err != nil {
			return nil, err
		}

		if booked {
			// Store booked slots to check for conflicts later
			bookedSlots[salesManagerID] = append(bookedSlots[salesManagerID], [2]time.Time{startDate, endDate})
		}
	}

	// Query again to find free slots (can optimize by combining with previous query)
	freeSlotQuery := `
      SELECT sales_manager_id, start_date
      FROM slots
      WHERE sales_manager_id = ANY($1)
      AND start_date::date = $2
      AND booked = FALSE
  `
	freeRows, err := db.Query(freeSlotQuery, "{"+managerIDs+"}", req.Date)
	if err != nil {
		return nil, err
	}
	defer freeRows.Close()

	// Check each free slot for conflicts and aggregate availability
	for freeRows.Next() {
		var salesManagerID int
		var startDate time.Time
		if err := freeRows.Scan(&salesManagerID, &startDate); err != nil {
			return nil, err
		}

		// Check if this slot conflicts with any booked slot for this manager
		conflict := false
		for _, bSlot := range bookedSlots[salesManagerID] {
			if startDate.Before(bSlot[1]) && startDate.Add(time.Hour).After(bSlot[0]) {
				// If the start time of the free slot is before the end of a booked slot
				// and the end time of the free slot is after the start of the booked slot, it conflicts.
				conflict = true
				break
			}
		}

		if !conflict {
			// Increment the count for this slot's start date
			formattedDate := startDate.Format("2006-01-02T15:04:05.000Z")
			slotAvailability[formattedDate]++
		}
	}

	// Convert the map to a slice of AvailableSlot
	var slots []AvailableSlot
	for startDate, count := range slotAvailability {
		slots = append(slots, AvailableSlot{
			AvailableCount: count,
			StartDate:      startDate,
		})
	}

	return slots, nil
}
