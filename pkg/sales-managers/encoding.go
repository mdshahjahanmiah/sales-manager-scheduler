package sales_managers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	eError "github.com/mdshahjahanmiah/explore-go/error"
)

// QueryRequest ...
type QueryRequest struct {
	Date     string   `json:"date"`
	Products []string `json:"products"`
	Language string   `json:"language"`
	Rating   string   `json:"rating"`
}

func decodeCalenderQuery(ctx context.Context, request *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(request.Body)

	var queryRequest QueryRequest
	err := decoder.Decode(&queryRequest)
	if err != nil {
		slog.Error("decode query request", "err", err)
		return nil, eError.NewServiceError(err, "decode decrypt request", "payload", http.StatusBadRequest)
	}
	if queryRequest.Date == "" {
		slog.Error("missing date")
		return nil, eError.NewServiceError(errors.New("date is empty"), "validation_error", "date", http.StatusBadRequest)
	}

	if queryRequest.Language == "" {
		slog.Error("missing language")
		return nil, eError.NewServiceError(errors.New("language is empty"), "validation_error", "language", http.StatusBadRequest)
	}

	if len(queryRequest.Products) == 0 {
		slog.Error("missing products")
		return nil, eError.NewServiceError(errors.New("product is empty"), "validation_error", "products", http.StatusBadRequest)
	}

	if queryRequest.Rating == "" {
		slog.Error("missing rating")
		return nil, eError.NewServiceError(errors.New("rating is empty"), "validation_error", "rating", http.StatusBadRequest)
	}

	return QueryRequest{
		Date:     queryRequest.Date,
		Products: queryRequest.Products,
		Language: queryRequest.Language,
		Rating:   queryRequest.Rating,
	}, nil
}
