package databaseconstants

import (
	"fmt"
	"strings"
)

type ImageType string

const (
	WorkerImage ImageType = "WORKER_IMAGE"
	LeaderImage ImageType = "LEADER_IMAGE"
)

func ParseImageType(s string) (ImageType, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case string(WorkerImage):
		return WorkerImage, nil
	case string(LeaderImage):
		return LeaderImage, nil
	default:
		return "", fmt.Errorf("invalid image type: %s", s)
	}
}
