package mocks

import "context"

type MockPricingRepository struct{}

func (m MockPricingRepository) GetPriceRule(t int, d int, c int64, golfers int) (float64, error) {
	return 1.25, nil
}

func (m MockPricingRepository) GetBasePrice(c int64) (float64, error) {
	return 12, nil
}

func (m MockPricingRepository) LookupBand(courseId string) (int64, error) {
	return 1, nil
}

func (m MockPricingRepository) UpdateBasePrice(ctx context.Context, courseId int64, year int64, price float64) error {
	return nil
}
