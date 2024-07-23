package environment

import (
	"github.com/jamiewhitney/fairways-core/pkg/cache"
	databasesql "github.com/jamiewhitney/fairways-core/pkg/database/mysql"
	"github.com/jamiewhitney/fairways-core/pkg/pubsub"
	"github.com/jamiewhitney/fairways-core/pkg/secretstore"
	booking_pb "github.com/jamiewhitney/fairways-core/protobufs/booking"
	pricing_pb "github.com/jamiewhitney/fairways-core/protobufs/pricing"
	tee_time_pb "github.com/jamiewhitney/fairways-core/protobufs/tee_time"
)

type Environment struct {
	secretstore   secretstore.SecretStore
	database      databasesql.DB
	cache         *cache.RedisRepository
	pubsub        *pubsub.Pubsub
	teeTimeClient tee_time_pb.TeeTimeServiceClient
	pricingClient pricing_pb.PricingServiceClient
	bookingClient booking_pb.BookingServiceClient

	Platform    string `env:"PLATFORM, default=local"`
	Environment string `env:"ENVIRONMENT, default=local"`
}

type Option func(env *Environment) *Environment

func New(env *Environment, opts ...Option) *Environment {
	for _, f := range opts {
		env = f(env)
	}

	return env
}

func WithSecretStore(sm secretstore.SecretStore) Option {
	return func(s *Environment) *Environment {
		s.secretstore = sm
		return s
	}
}

func WithDatabase(database *databasesql.DB) Option {
	return func(s *Environment) *Environment {
		s.database = *database
		return s
	}
}

func WithCache(cache cache.RedisRepository) Option {
	return func(s *Environment) *Environment {
		s.cache = &cache
		return s
	}
}

func WithTeeTimeServiceClient(serviceClient tee_time_pb.TeeTimeServiceClient) Option {
	return func(s *Environment) *Environment {
		s.teeTimeClient = serviceClient
		return s
	}
}

func WithPubsubConfig(pubsub *pubsub.Pubsub) Option {
	return func(s *Environment) *Environment {
		s.pubsub = pubsub
		return s
	}
}

func WithPricingServiceClient(serviceClient pricing_pb.PricingServiceClient) Option {
	return func(s *Environment) *Environment {
		s.pricingClient = serviceClient
		return s
	}
}

func WithBookingServiceClient(serviceClient booking_pb.BookingServiceClient) Option {
	return func(s *Environment) *Environment {
		s.bookingClient = serviceClient
		return s
	}
}

func (e *Environment) PlatformC() string {
	return e.Platform
}

func (e *Environment) EnvironmentC() string {
	return e.Environment
}

func (e *Environment) Database() *databasesql.DB {
	return &e.database
}

func (e *Environment) Cache() *cache.RedisRepository {
	return e.cache
}

func (e *Environment) Pubsub() pubsub.Pubsub {
	return *e.pubsub
}

func (e *Environment) BookingClient() booking_pb.BookingServiceClient {
	return e.bookingClient

}

func (e *Environment) PricingClient() pricing_pb.PricingServiceClient {
	return e.pricingClient

}

func (e *Environment) TeeTimeClient() tee_time_pb.TeeTimeServiceClient {
	return e.teeTimeClient
}

func (e *Environment) GetEnvironment() string {
	return e.Environment
}

func (e *Environment) GetPlatform() string {
	return e.Platform
}
