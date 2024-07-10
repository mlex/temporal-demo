package helloworld

import (
	"time"

	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
	"go.temporal.io/sdk/workflow"
)

func MyWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	now := workflow.Now(ctx)
	logger.Info("It's now " + now.String())
	logger.Info("Sleeping until " + now.Add(30*time.Second).String())

	workflow.Sleep(ctx, 20*time.Second)

	now = workflow.Now(ctx)
	logger.Info("Done with sleeping. It's now " + now.String())

	return "workflow-result", nil
}
