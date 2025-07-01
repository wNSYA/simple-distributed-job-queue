package _interface

import (
	"context"
	"jobqueue/entity"
)

type JobService interface {
	Enqueue(ctx context.Context, taskName string) (*entity.Job, error)
	GetAllJobs(ctx context.Context) (output []*entity.Job, err error)
	FindByID(ctx context.Context, id string) (*entity.Job, error)
}

type JobRepository interface {
	Save(ctx context.Context, job *entity.Job) error
	FindByID(ctx context.Context, id string) (*entity.Job, error)
	FindAll(ctx context.Context) ([]*entity.Job, error)
}
