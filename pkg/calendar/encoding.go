package calendar

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	eError "github.com/mdshahjahanmiah/explore-go/error"
)

type queryRequest struct {
	Date     string   `json:"date"`
	Products []string `json:"products"`
	Language string   `json:"language"`
	Rating   string   `json:"rating"`
}

func decodeCalenderQuery(_ context.Context, r *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)

	var request queryRequest
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("error during decoding query request", "err", err)
		return nil, eError.NewServiceError(err, "decode query request", "payload", http.StatusBadRequest)
	}

	fieldChecks := map[string]string{
		"date":     request.Date,
		"language": request.Language,
		"rating":   request.Rating,
	}

	for field, value := range fieldChecks {
		if value == "" {
			return nil, eError.NewServiceError(errors.New("field is empty"), "validation_error", field, http.StatusBadRequest)
		}
	}

	if len(request.Products) == 0 {
		return nil, eError.NewServiceError(errors.New("filed is empty"), "validation_error", "products", http.StatusBadRequest)
	}

	return queryRequest{
		Date:     request.Date,
		Products: request.Products,
		Language: request.Language,
		Rating:   request.Rating,
	}, nil
}
