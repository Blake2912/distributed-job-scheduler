package container

import (
	"log"
	"time"

	"example.com/httpclient"
	spawnworkersHandler "example.com/leader/internal/api/handlers/spawn_workers"
	spawnworkers "example.com/leader/internal/services/spawn_workers"

	"example.com/pod_library/client"
)

type Container struct {
	SpawnWorkersHandler *spawnworkersHandler.SpawnWorkersHandler
}

func BuildContainer() *Container {

	httpClient := httpclient.New(120 * time.Minute)
	k8sClient, err := client.New()
	if err != nil {
		log.Fatal(err)
	}

	// Build services
	spawnWorkerService := spawnworkers.NewSpawnWorkerService(*httpClient, *k8sClient)

	// Build Handlers
	spawnWorkersHandler := spawnworkersHandler.NewSpawnWokersHandler(spawnWorkerService)

	return &Container{
		SpawnWorkersHandler: spawnWorkersHandler,
	}
}
