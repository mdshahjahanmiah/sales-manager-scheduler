package calendar

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDecodeCalenderQuery_Success(t *testing.T) {
	body := `{
		"date": "2024-05-03",
		"products": ["SolarPanels", "Heatpumps"],
		"language": "German",
		"rating": "Gold"
	}`
	req := httptest.NewRequest(http.MethodPost, "/calendar/query", bytes.NewBufferString(body))
	result, err := decodeCalenderQuery(context.Background(), req)

	// Check if the function returned an error when it shouldn't have
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Define the expected result
	expected := queryRequest{
		Date:     "2024-05-03",
		Products: []string{"SolarPanels", "Heatpumps"},
		Language: "German",
		Rating:   "Gold",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestDecodeCalenderQuery_Fail(t *testing.T) {
	// Define the input body that represents an invalid request (missing required fields)
	body := `{
		"products": ["SolarPanels", "Heatpumps"],
		"language": "German",
		"rating": "Gold"
	}`

	req := httptest.NewRequest(http.MethodPost, "/calendar/query", bytes.NewBufferString(body))
	result, err := decodeCalenderQuery(context.Background(), req)

	// Check if the function returned an error when it should have
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}
