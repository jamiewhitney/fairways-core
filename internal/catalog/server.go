package catalog

import (
	"context"
	"github.com/jamiewhitney/fairways-core/internal/catalog/services"
	"github.com/jamiewhitney/fairways-core/pkg/environment"
	"github.com/jamiewhitney/fairways-core/protobufs/catalog"
)

type Server struct {
	catalog.UnimplementedCatalogServiceServer

	catalogService *services.CatalogService
}

func New(env *environment.Environment, config *Config) Server {
	return Server{
		catalogService: services.NewCatalogService(env.Database(), env.Cache(), env.Pubsub()),
	}
}

func (s *Server) CourseExists(ctx context.Context, in *catalog.CourseExistsRequest) (*catalog.CourseExistsResponse, error) {
	_, err := s.catalogService.GetCourse(ctx, in.CourseId)
	if err != nil {
		return nil, err
	}
	return &catalog.CourseExistsResponse{Exists: true}, nil
}

func (s *Server) GetCourses(ctx context.Context, in *catalog.GetCoursesRequest) (*catalog.GetCoursesResponse, error) {
	return s.catalogService.ListCourses(ctx, in.Offset, in.Limit)

}
func (s *Server) GetCourse(ctx context.Context, in *catalog.GetCourseRequest) (*catalog.Course, error) {
	return s.catalogService.GetCourse(ctx, in.CourseId)
}

func (s *Server) CreateCourse(ctx context.Context, in *catalog.CreateCourseRequest) (*catalog.Course, error) {
	return s.catalogService.CreateCourse(ctx, in)
}
