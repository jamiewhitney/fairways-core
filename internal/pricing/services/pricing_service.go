package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jamiewhitney/fairways-core/internal/pricing/repository"
	"github.com/jamiewhitney/fairways-core/pkg/cache"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	tee_time_pb "github.com/jamiewhitney/fairways-core/protobufs/tee_time"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type PricingService struct {
	pricingRepository repository.PricingRepository
	teeTimeService    tee_time_pb.TeeTimeServiceClient
	cache             *cache.RedisRepository
}

func NewPricingService(pricingRepository repository.PricingRepository, cache *cache.RedisRepository) *PricingService {
	return &PricingService{
		pricingRepository: pricingRepository,
		cache:             cache,
	}
}

func (this *PricingService) GetPrice(ctx context.Context, t string, courseId string, golfers int) (map[string]interface{}, error) {
	logger := logging.FromContext(ctx)
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return nil, err
	}

	_, week := parsedTime.ISOWeek()

	band, err := this.LookupBand(courseId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error(err)
			return nil, errors.New("course is not associated with a band")
		}
		logger.Error(err)
		return nil, err
	}

	modifier, err := this.pricingRepository.GetPriceRule(week, parsedTime.Hour(), band, golfers)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithFields(logrus.Fields{
				"week": week,
				"hour": parsedTime.Hour(),
				"band": band,
			}).Debug("no pricing rules")

		} else {
			logger.Error(err)
			return nil, err
		}
	}

	cid, err := strconv.ParseInt(courseId, 10, 64)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	basePrice, err := this.GetBasePrice(cid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error(err)
			return nil, errors.New(fmt.Sprintf("course %d does not havea base price", cid))
		}
		logger.Error(err)
		return nil, err
	}

	var result float64

	if modifier == 0 {
		result = basePrice
	} else {
		result = basePrice * modifier
	}
	fmt.Println(modifier)
	fmt.Println(result)
	fmt.Println(basePrice)

	logger.Debugf("result: %v base_price: %v modifier: %v", result, basePrice, modifier)
	if result < basePrice {
		return map[string]interface{}{
			"base_price": basePrice,
			"price":      result,
			"discounted": true,
		}, nil
	}

	return map[string]interface{}{
		"base_price": basePrice,
		"price":      result,
		"discounted": false,
	}, nil
}

func (this *PricingService) LookupBand(courseId string) (int64, error) {
	//result, err := this.pricingRepository.LookupBand(courseId)
	//if err != nil {
	//	return 0, err
	//}
	//return result, nil
	return 1, nil
}

func (this *PricingService) GetBasePrice(c int64) (float64, error) {
	result, err := this.pricingRepository.GetBasePrice(c)
	if err != nil {
		return 0, err
	}
	return result, nil
}
