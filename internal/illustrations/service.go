package illustrations

import (
	"context"
	"fmt"

	repo "github.com/fillipgms/portfolio-api/internal/adapters/postgresql/sqlc"
)

type Service interface{
	CreateIllustration(ctx context.Context, tempIllustration repo.CreateIllustrationParams) (repo.Illustration, error)
	ListIllustrations(ctx context.Context, limit int32, offset int32) ([]repo.Illustration, error);
	FindIllustrationsCount(ctx context.Context) (int64, error);
	FindIllustrationById(ctx context.Context, id int64) (repo.Illustration, error)
}

type svc struct {
	repo repo.Querier 
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) CreateIllustration(ctx context.Context, tempIllustration repo.CreateIllustrationParams) (repo.Illustration, error) {
	if tempIllustration.Title == "" {
		return repo.Illustration{}, fmt.Errorf("Title is Required")
	}

	if tempIllustration.Description == "" {
		return repo.Illustration{}, fmt.Errorf("Description is Required")
	}

	if tempIllustration.Imageurl == "" {
		return repo.Illustration{}, fmt.Errorf("Image Url is Required")
	}
	
	return s.repo.CreateIllustration(ctx, tempIllustration)
}

func (s *svc) ListIllustrations(ctx context.Context, limit int32, offset int32) ([]repo.Illustration, error) {
	params := repo.ListIllustrationsParams{
		Limit:  limit,
		Offset: offset,
	}

	return s.repo.ListIllustrations(ctx, params)
}

func (s *svc) FindIllustrationsCount(ctx context.Context) (int64, error) {
	return s.repo.FindIllustrationsCount(ctx)
}

func (s *svc) FindIllustrationById(ctx context.Context, id int64) (repo.Illustration, error) {
	return s.repo.FindIllustrationById(ctx, id)
}



