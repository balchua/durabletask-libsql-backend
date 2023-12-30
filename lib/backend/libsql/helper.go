package libsql

import (
	"github.com/microsoft/durabletask-go/api"
	"github.com/microsoft/durabletask-go/backend"
)

var orchestrationStatus = map[string]api.OrchestrationStatus{
	"ORCHESTRATION_STATUS_RUNNING":          api.RUNTIME_STATUS_RUNNING,
	"ORCHESTRATION_STATUS_COMPLETED":        api.RUNTIME_STATUS_COMPLETED,
	"ORCHESTRATION_STATUS_CONTINUED_AS_NEW": api.RUNTIME_STATUS_CONTINUED_AS_NEW,
	"ORCHESTRATION_STATUS_FAILED":           api.RUNTIME_STATUS_FAILED,
	"ORCHESTRATION_STATUS_CANCELED":         api.RUNTIME_STATUS_CANCELED,
	"ORCHESTRATION_STATUS_TERMINATED":       api.RUNTIME_STATUS_TERMINATED,
	"ORCHESTRATION_STATUS_PENDING":          api.RUNTIME_STATUS_PENDING,
	"ORCHESTRATION_STATUS_SUSPENDED":        api.RUNTIME_STATUS_SUSPENDED,
}

func toRuntimeStatusString(status api.OrchestrationStatus) string {
	return status.String()
}

func fromRuntimeStatusString(status string) api.OrchestrationStatus {
	return orchestrationStatus[status]
}

func toExecutionStartedEvent(historyEvent *backend.HistoryEvent) *executionStartedEvent {
	startEvent := historyEvent.GetExecutionStarted()
	return &executionStartedEvent{
		name:        startEvent.Name,
		version:     startEvent.Version.GetValue(),
		instanceId:  startEvent.OrchestrationInstance.InstanceId,
		executionId: startEvent.OrchestrationInstance.ExecutionId.GetValue(),
		input:       startEvent.Input.GetValue(),
	}
}

type executionStartedEvent struct {
	name        string
	version     string
	instanceId  string
	executionId string
	input       string
}
