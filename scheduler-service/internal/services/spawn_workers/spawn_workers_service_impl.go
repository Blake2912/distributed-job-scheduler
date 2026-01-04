package spawnworkers

import (
	"context"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/httpclient"
	podLibraryClient "github.com/Blake2912/distributed-job-scheduler/scheduler-service/pod_library/client"
)

type spawnWorkersService struct {
	httpClient httpclient.Client
	k8sClient  podLibraryClient.K8sClient
}

func NewSpawnWorkerService(
	httpClient httpclient.Client,
	k8sClient podLibraryClient.K8sClient,
) SpawnWorkersService {
	return &spawnWorkersService{
		httpClient: httpClient,
		k8sClient:  k8sClient,
	}
}

func (s *spawnWorkersService) SpawnWorkers(ctx context.Context, noOfWorkers int) error {

	// Fetch the worker image from DAL

	// Call the deployment function from the pod library

	return nil
}
