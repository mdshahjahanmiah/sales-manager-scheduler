package calendar

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makePostCalenderQueryEndpoint(ms Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(queryRequest)

		result, err := ms.AvailableSlots(req)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}
