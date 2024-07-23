package services

import (
	"context"
	"fmt"
	"github.com/jamiewhitney/fairways-core/internal/catalog/repository"
	"github.com/jamiewhitney/fairways-core/pkg/cache"
	databasesql "github.com/jamiewhitney/fairways-core/pkg/database/mysql"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	"github.com/jamiewhitney/fairways-core/pkg/pubsub"
	"github.com/jamiewhitney/fairways-core/protobufs/catalog"
)

type CatalogService struct {
	catalogRepository repository.Querier
	cache             *cache.RedisRepository
	Pubsub            pubsub.Pubsub
	db                *databasesql.DB
	queries           *repository.Queries
}

func NewCatalogService(db *databasesql.DB, cache *cache.RedisRepository, pubsub pubsub.Pubsub) *CatalogService {
	return &CatalogService{
		catalogRepository: repository.New(db.Pool),
		cache:             cache,
		Pubsub:            pubsub,
		db:                db,
	}
}

func (cs *CatalogService) GetCourse(ctx context.Context, courseId int64) (*catalog.Course, error) {
	logger := logging.FromContext(ctx)
	result, err := cs.catalogRepository.GetACourse(ctx, courseId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &catalog.Course{
		Id:       result.ID,
		Name:     result.Name,
		Town:     result.City,
		County:   result.State,
		Postcode: result.PostalCode,
	}, nil
}

func (cs *CatalogService) ListCourses(ctx context.Context, requestedOffset int64, requestedLimit int64) (*catalog.GetCoursesResponse, error) {
	logger := logging.FromContext(ctx)

	if requestedLimit > 50 {
		return nil, fmt.Errorf("requested result cannot be greater than 50")
	}

	if requestedLimit == 0 {
		requestedLimit = 20
	}

	results, err := cs.catalogRepository.ListCourses(ctx, repository.ListCoursesParams{
		Limit:  int32(requestedLimit),
		Offset: int32(requestedOffset),
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var courses []*catalog.Course
	for _, result := range results {
		courses = append(courses, &catalog.Course{
			Id:       result.ID,
			Name:     result.Name,
			Town:     result.City,
			County:   result.State,
			Postcode: result.PostalCode,
		})
	}

	return &catalog.GetCoursesResponse{
		Courses: courses,
		Limit:   requestedLimit,
		Offset:  requestedOffset,
	}, nil
}

func (cs *CatalogService) CreateCourse(ctx context.Context, in *catalog.CreateCourseRequest) (*catalog.Course, error) {
	logger := logging.FromContext(ctx)

	tx, err := cs.db.Pool.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer tx.Rollback()

	queries := cs.queries.WithTx(tx)

	courseId, err := queries.CreateCourse(ctx, repository.CreateCourseParams{
		Name:       in.Name,
		City:       in.Town,
		State:      in.County,
		PostalCode: in.Postcode,
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	course, err := queries.GetACourse(ctx, courseId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if err := cs.Pubsub.PublishMessage(ctx, "catalog-update", fmt.Sprintf("New course created: %s", course.Name)); err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.Error(err)
		return nil, err
	}
	return &catalog.Course{
		Id:       course.ID,
		Name:     course.Name,
		Town:     course.City,
		County:   course.State,
		Postcode: course.PostalCode,
		Live:     course.Live,
	}, nil
}
