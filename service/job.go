package service

import (
	"context"
	"jobqueue/entity"
	_interface "jobqueue/interface"
)

type jobService struct {
	jobRepo _interface.JobRepository
}

// Initiator ...
type Initiator func(s *jobService) *jobService

func (q jobService) GetAllJobs(ctx context.Context) (output entity.Job, err error) {
	output = entity.Job{}
	return output, nil
}

func (q jobService) Enqueue(ctx context.Context, taskName string) (string, error) {
	retval := "ok"
	return retval, nil
}

// NewJobService ...
func NewJobService() Initiator {
	return func(s *jobService) *jobService {
		return s
	}
}

// SetJobRepository ...
func (i Initiator) SetJobRepository(jobRepository _interface.JobRepository) Initiator {
	return func(s *jobService) *jobService {
		i(s).jobRepo = jobRepository
		return s
	}
}

// Build ...
func (i Initiator) Build() _interface.JobService {
	return i(&jobService{})
}
