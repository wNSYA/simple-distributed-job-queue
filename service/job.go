package service

import (
	"context"
	"fmt"
	"jobqueue/entity"
	_interface "jobqueue/interface"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type jobService struct {
	jobRepo _interface.JobRepository
	mu      sync.Mutex
}

// Initiator ...
type Initiator func(s *jobService) *jobService

func (q *jobService) GetAllJobs(ctx context.Context) (output []*entity.Job, err error) {
	fmt.Println("from service/job GetAllJobs")
	return q.jobRepo.FindAll(ctx)
}

func (q *jobService) FindByID(ctx context.Context, id string) (*entity.Job, error) {
	fmt.Println("from service/job FindByID")
	return q.jobRepo.FindByID(ctx, id)
}

func (q *jobService) Enqueue(ctx context.Context, taskName string) (*entity.Job, error) {
	fmt.Println("from service/job Enqueue")

	jobs, _ := q.jobRepo.FindAll(ctx)

	hasRunning := false
	for _, job := range jobs {
		if job.Status == "Running" {
			hasRunning = true
			break
		}
	}

	job := &entity.Job{
		ID:       uuid.NewV4().String(),
		Task:     taskName,
		Status:   "Pending",
		Attempts: 0,
	}

	if !hasRunning {
		job.Status = "Running"
		go q.processJob(ctx, job)
	}

	err := q.jobRepo.Save(ctx, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (s *jobService) processJob(ctx context.Context, job *entity.Job) {
	time.Sleep(5 * time.Second)

	s.mu.Lock()
	defer s.mu.Unlock()

	job.Attempts++

	if job.Task == "unstable-job" {
		if job.Attempts < 3 {
			job.Status = "Failed"
			fmt.Printf("Job %s failed on attempt %d\n", job.ID, job.Attempts)

			// Re-enqueue the job as pending
			job.Status = "Pending"

			// Schedule the job again
			go s.promoteNextPendingJob(ctx)
			return
		}
		// On 3rd attempt, complete the job
		job.Status = "Completed"
		fmt.Printf("Job %s completed\n", job.ID)
	} else {
		// Normal job completes
		job.Status = "Completed"
		fmt.Printf("Job %s completed\n", job.ID)
	}

	// Trigger next job in queue
	go s.promoteNextPendingJob(ctx)

}

func (s *jobService) promoteNextPendingJob(ctx context.Context) {
	jobs, _ := s.jobRepo.FindAll(ctx)
	for _, job := range jobs {
		if job.Status == "Pending" {
			fmt.Printf("Promoting job %s to running\n", job.ID)

			job.Status = "Running"
			go s.processJob(ctx, job)
			break
		}
	}
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
