package sales_managers

import (
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/mdshahjahanmiah/explore-go/error"
	"github.com/mdshahjahanmiah/explore-go/http"
)

func MakeHandler(ms Service) http.Endpoint {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(error.EncodeError),
	}

	getMidpointBandHandler := kithttp.NewServer(
		makePostCalenderQueryEndpoint(ms),
		decodeCalenderQuery,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	r := chi.NewRouter()

	r.Method("POST", "/calendar/query", getMidpointBandHandler)

	return http.Endpoint{Pattern: "/*", Handler: r}
}