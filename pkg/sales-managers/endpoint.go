package sales_managers

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makePostCalenderQueryEndpoint(ms Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		queryRequest := request.(QueryRequest)

		result, err := ms.AvailableSlots(queryRequest)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}
