package mutation

import (
	"context"
	"fmt"
	_dataloader "jobqueue/delivery/graphql/dataloader"
	"jobqueue/delivery/graphql/resolver"
	_interface "jobqueue/interface"

	"jobqueue/entity"
)

type JobMutation struct {
	jobService _interface.JobService
	dataloader *_dataloader.GeneralDataloader
}

func (q JobMutation) Enqueue(ctx context.Context, args entity.Job) (*resolver.JobResolver, error) {
	fmt.Println("from mutation/job")
	job, err := q.jobService.Enqueue(ctx, args.Task)
	if err != nil {
		fmt.Printf("Failed to enqueue job: %v\n", err)
		return nil, err
	}
	return &resolver.JobResolver{
		Data:       *job,
		JobService: q.jobService,
		Dataloader: q.dataloader,
	}, nil
}

// NewJobMutation to create new instance
func NewJobMutation(jobService _interface.JobService, dataloader *_dataloader.GeneralDataloader) JobMutation {
	return JobMutation{
		jobService: jobService,
		dataloader: dataloader,
	}
}
