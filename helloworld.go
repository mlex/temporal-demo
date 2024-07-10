package helloworld

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"go.temporal.io/sdk/activity"
	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a Hello World workflow definition.
func MyWorkflow(ctx workflow.Context, name string) (string, error) {
	var signal interface{}
	var signalReceived bool
	signalChan := workflow.GetSignalChannel(ctx, "my-signal")
	workflow.Go(ctx, func(ctx workflow.Context) {
		for {
			selector := workflow.NewSelector(ctx)
			selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
				c.Receive(ctx, &signal)
				signalReceived = true
			})
			selector.Select(ctx)
		}
	})

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	now := workflow.Now(ctx)
	logger.Info("It's now " + now.String())
	logger.Info("Sleeping until " + now.Add(30*time.Second).String())

	ok, _ := workflow.AwaitWithTimeout(ctx, 30*time.Second, func() bool {
		return signalReceived
	})

	if ok {
		msg, _ := json.Marshal(signal)
		logger.Info("Received signal: " + string(msg))
	} else {
		now = workflow.Now(ctx)
		logger.Info("Waited too long - it's now " + now.String())
	}

	var result string
	err := workflow.ExecuteActivity(ctx, MyActivity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("HelloWorld workflow completed.", "result", result)

	return result, nil
}

func MyActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("MyActivity", "name", name)
	data, err := os.ReadFile("/tmp/dat")
	return string(data), err
}
