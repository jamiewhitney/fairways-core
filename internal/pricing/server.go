package pricing

import (
	"context"
	"fmt"
	"github.com/jamiewhitney/fairways-core/internal/pricing/repository"
	"github.com/jamiewhitney/fairways-core/internal/pricing/services"
	"github.com/jamiewhitney/fairways-core/pkg/environment"
	pricing_pb "github.com/jamiewhitney/fairways-core/protobufs/pricing"
	"io"
	"strconv"
)

type Server struct {
	pricing_pb.UnimplementedPricingServiceServer

	pricingService *services.PricingService
}

func New(env *environment.Environment, config *Config) Server {
	return Server{
		pricingService: services.NewPricingService(repository.NewMySQLPricingRepository(env.Database()), env.Cache()),
	}
}

func (s *Server) GetPrice(ctx context.Context, in *pricing_pb.GetPriceRequest) (*pricing_pb.GetPriceResponse, error) {
	result, err := s.pricingService.GetPrice(ctx, in.Datetime, strconv.FormatInt(in.CourseId, 10), int(in.Golfers))
	if err != nil {
		return nil, err
	}

	return &pricing_pb.GetPriceResponse{
		Price:    result["price"].(float64),
		CourseId: strconv.FormatInt(in.CourseId, 10),
		Datetime: in.Datetime,
	}, nil
}

func (s *Server) GetPriceStream(stream pricing_pb.PricingService_GetPriceStreamServer) error {
	if stream == nil {
		return fmt.Errorf("stream is nil")
	}
	if s == nil {
		return fmt.Errorf("RPC receiver is nil")
	}
	if s.pricingService == nil {
		return fmt.Errorf("pricingController is nil")
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if req == nil {
			return fmt.Errorf("received nil request")
		}

		result, err := s.pricingService.GetPrice(stream.Context(),
			req.Datetime, strconv.FormatInt(req.CourseId, 10), int(req.Golfers))
		if err != nil {
			return err
		}

		response := &pricing_pb.GetPriceResponse{
			Price:         result["price"].(float64),
			CourseId:      strconv.FormatInt(req.CourseId, 10),
			Datetime:      req.Datetime,
			OriginalPrice: result["base_price"].(float64),
			Discounted:    result["discounted"].(bool),
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}
