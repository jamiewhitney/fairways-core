package services

import (
	"context"
	mocks "github.com/jamiewhitney/fairways-core/testing/mocks/repositories"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

var (
	ps *PricingService
)

func init() {
	logger := logrus.New()
	logger.Out = io.Discard

	h := mocks.MockPricingRepository{}
	ps = NewPricingService(h, nil)
}
func TestDefaultPricingController_GetPrice(t *testing.T) {
	tests := []struct {
		name     string
		time     string
		courseId string
		golfers  int
	}{
		{
			name:     "TestDefaultPricingController_GetPrice",
			time:     "2020-01-01T09:12:00Z",
			courseId: "10046",
			golfers:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ps.GetPrice(context.Background(), tt.time, tt.courseId, tt.golfers)
			assert.NotNil(t, result)
			assert.NoError(t, err)
		})
	}
}
