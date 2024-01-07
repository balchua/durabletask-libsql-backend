package simple

import (
	"log/slog"

	"github.com/microsoft/durabletask-go/task"
)

type employee struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

var employees = []employee{
	{"1", "John", "Doe", 30},
	{"2", "Jane", "Doe", 25},
	{"3", "Steven", "Smith", 40},
	{"4", "Lily", "Smith", 35},
}

type SimpleWorkflow struct {
}

func NewSimpleWorkflow() *SimpleWorkflow {
	return &SimpleWorkflow{}
}

func (s *SimpleWorkflow) GetWorkflow() func(ctx *task.OrchestrationContext) (any, error) {
	return s.simpleOrchestration
}
func (s *SimpleWorkflow) GetTasks() map[string]func(ctx task.ActivityContext) (any, error) {

	return map[string]func(ctx task.ActivityContext) (any, error){
		"GetEmployeeDetailById": s.getEmployeeDetailById,
		"Task2":                 s.task2,
	}
}
func (s *SimpleWorkflow) simpleOrchestration(ctx *task.OrchestrationContext) (any, error) {
	var input int
	ctx.GetInput(&input)
	slog.Info("input", "value", input)

	var john employee
	if err := ctx.CallActivity("GetEmployeeDetailById", task.WithActivityInput("1")).Await(&john); err != nil {
		return nil, err
	}
	slog.Info("john's reponse", "value", john)
	var lily employee
	if err := ctx.CallActivity("GetEmployeeDetailById", task.WithActivityInput("4")).Await(&lily); err != nil {
		return nil, err
	}
	slog.Info("lily's reponse", "value", lily)

	var jane employee
	if err := ctx.CallActivity("GetEmployeeDetailById", task.WithActivityInput("2")).Await(&jane); err != nil {
		return nil, err
	}
	slog.Info("jane's reponse", "value", jane)

	var steven employee
	if err := ctx.CallActivity("GetEmployeeDetailById", task.WithActivityInput("3")).Await(&steven); err != nil {
		return nil, err
	}
	slog.Info("steven's reponse", "value", steven)

	var task2Response string
	if err := ctx.CallActivity("Task2", task.WithActivityInput("3")).Await(&task2Response); err != nil {
		return nil, err
	}

	slog.Info("task 2's reponse", "value", task2Response)

	return []employee{john, lily, jane, steven}, nil
}
