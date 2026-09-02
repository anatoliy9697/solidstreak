package schedulertask

type TaskType string

const (
	TaskTypeProcessExpiredInvoice TaskType = "process_expired_invoice"
)

var TaskTypeMapping = map[string]TaskType{
	string(TaskTypeProcessExpiredInvoice): TaskTypeProcessExpiredInvoice,
}

type Task interface {
	Type() TaskType
}
