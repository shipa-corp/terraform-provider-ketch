package client

import (
	"context"
	"log"

	"github.com/brunoa19/ketch-terraform-provider/client/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Container struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Command []string `json:"command"`
}

type Policy struct {
	RestartPolicy string `json:"restart_policy,omitempty"`
}

type Job struct {
	Version      string       `json:"version,omitempty"`
	Type         string       `json:"type"`
	Name         string       `json:"name"`
	Framework    string       `json:"framework"`
	Description  string       `json:"description,omitempty"`
	Parallelism  int64        `json:"parallelism,omitempty"`
	Completions  int64        `json:"completions,omitempty"`
	Suspend      bool         `json:"suspend,omitempty"`
	BackoffLimit int64        `json:"backoff_limit,omitempty"`
	Containers   []*Container `json:"containers,omitempty"`
	Policy       *Policy      `json:"policy,omitempty"`
}

func NewJob(input *v1beta1.Job) *Job {
	var containers []*Container
	for _, container := range input.Spec.Containers {
		containers = append(containers, &Container{
			Name:    container.Name,
			Image:   container.Image,
			Command: container.Command,
		})
	}

	return &Job{
		Version:      input.Spec.Version,
		Type:         input.Spec.Type,
		Name:         input.Spec.Name,
		Framework:    input.Spec.Framework,
		Description:  input.Spec.Description,
		Parallelism:  int64(input.Spec.Parallelism),
		Completions:  int64(input.Spec.Completions),
		Suspend:      input.Spec.Suspend,
		BackoffLimit: int64(input.Spec.BackoffLimit),
		Containers:   containers,
		Policy: &Policy{
			RestartPolicy: string(input.Spec.Policy.RestartPolicy),
		},
	}
}

func (j *Job) convertToKetchJob() *v1beta1.Job {
	var containers []v1beta1.Container
	for _, container := range j.Containers {
		containers = append(containers, v1beta1.Container{
			Name:    container.Name,
			Image:   container.Image,
			Command: container.Command,
		})
	}

	var restartPolicy v1beta1.RestartPolicy
	if j.Policy != nil {
		restartPolicy = v1beta1.RestartPolicy(j.Policy.RestartPolicy)
		if restartPolicy != v1beta1.Never && restartPolicy != v1beta1.OnFailure {
			restartPolicy = ""
		}
	}

	spec := v1beta1.JobSpec{
		Version:      j.Version,
		Type:         j.Type,
		Name:         j.Name,
		Framework:    j.Framework,
		Description:  j.Description,
		Parallelism:  int(j.Parallelism),
		Completions:  int(j.Completions),
		Suspend:      j.Suspend,
		BackoffLimit: int(j.BackoffLimit),
		Containers:   containers,
		Policy: v1beta1.Policy{
			RestartPolicy: restartPolicy,
		},
	}

	setJobSpecDefaults(&spec)

	return &v1beta1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: j.Name,
		},
		Spec: spec,
	}
}

const (
	defaultJobVersion       = "v1"
	defaultJobParallelism   = 1
	defaultJobCompletions   = 1
	defaultJobBackoffLimit  = 6
	defaultJobRestartPolicy = "Never"
)

func setJobSpecDefaults(jobSpec *v1beta1.JobSpec) {
	jobSpec.Type = "Job"
	if jobSpec.Version == "" {
		jobSpec.Version = defaultJobVersion
	}
	if jobSpec.Parallelism == 0 {
		jobSpec.Parallelism = defaultJobParallelism
	}
	if jobSpec.Completions == 0 {
		jobSpec.Completions = defaultJobCompletions
	}
	if jobSpec.BackoffLimit == 0 {
		jobSpec.BackoffLimit = defaultJobBackoffLimit
	}
	if jobSpec.Policy.RestartPolicy == "" {
		jobSpec.Policy.RestartPolicy = defaultJobRestartPolicy
	}
}

func (c *Client) getJob(ctx context.Context, name string) (*v1beta1.Job, error) {
	job := &v1beta1.Job{}
	err := c.kube.Get(ctx, types.NamespacedName{Name: name}, job)
	if err != nil {
		log.Println("ERR:", err.Error())
		return nil, err
	}

	return job, nil
}

func (c *Client) CreateJob(ctx context.Context, input *Job) error {
	job := input.convertToKetchJob()
	return c.kube.Create(ctx, job)
}

func (c *Client) UpdateJob(ctx context.Context, input *Job) error {
	job, err := c.getJob(ctx, input.Name)
	if err != nil {
		return err
	}

	updates := input.convertToKetchJob()
	job.Spec = updates.Spec
	return c.kube.Update(ctx, job)
}

func (c *Client) GetJob(ctx context.Context, name string) (*Job, error) {
	job, err := c.getJob(ctx, name)
	if err != nil {
		return nil, err
	}

	return NewJob(job), nil
}

func (c *Client) DeleteJob(ctx context.Context, name string) error {
	job, err := c.getJob(ctx, name)
	if err != nil {
		return err
	}

	err = c.kube.Delete(ctx, job)
	if err != nil {
		return err
	}

	return nil
}
