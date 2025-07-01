package query

import (
	"context"
	"fmt"
	_dataloader "jobqueue/delivery/graphql/dataloader"
	"jobqueue/delivery/graphql/resolver"
	"jobqueue/entity"
	_interface "jobqueue/interface"
)

type JobQuery struct {
	jobService _interface.JobService
	dataloader *_dataloader.GeneralDataloader
}

func (q JobQuery) Jobs(ctx context.Context) ([]resolver.JobResolver, error) {
	fmt.Println("from query/job Jobs")

	jobs, err := q.jobService.GetAllJobs(ctx)
	if err != nil {
		return nil, err
	}

	resolvers := make([]resolver.JobResolver, 0, len(jobs))
	for _, job := range jobs {
		resolvers = append(resolvers, resolver.JobResolver{
			Data:       *job,
			JobService: q.jobService,
			Dataloader: q.dataloader,
		})
	}

	return resolvers, nil
}

func (q JobQuery) Job(ctx context.Context, args struct{ ID string }) (*resolver.JobResolver, error) {
	fmt.Println("from query/job JobID")

	job, err := q.jobService.FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &resolver.JobResolver{
		Data:       *job,
		JobService: q.jobService,
		Dataloader: q.dataloader,
	}, nil
}

func (q JobQuery) JobStatus(ctx context.Context) (resolver.JobStatusResolver, error) {
	fmt.Println("from query/job status")

	jobs, err := q.jobService.GetAllJobs(ctx)
	if err != nil {
		return resolver.JobStatusResolver{}, err
	}

	statusCount := entity.JobStatus{}
	for _, job := range jobs {
		switch job.Status {
		case "Completed":
			statusCount.Completed++
		case "Pending":
			statusCount.Pending++
		case "Running":
			statusCount.Running++
		case "Failed":
			statusCount.Failed++
		}
	}

	return resolver.JobStatusResolver{
		Data:       statusCount,
		JobService: q.jobService,
		Dataloader: q.dataloader,
	}, nil
}

func NewJobQuery(jobService _interface.JobService,
	dataloader *_dataloader.GeneralDataloader) JobQuery {
	return JobQuery{
		jobService: jobService,
		dataloader: dataloader,
	}
}
