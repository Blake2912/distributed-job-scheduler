package helpers

import (
	"fmt"

	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
)

func BuildHealthCheckKey(workerId string, jobExecutionId string) string {
	return fmt.Sprintf("%s#%s#%s", database_constants.HEALTH_CHECK_KEY_IDENTIFIER, workerId, jobExecutionId)
}
