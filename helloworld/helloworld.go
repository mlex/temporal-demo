package helloworld

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"go.temporal.io/sdk/activity"
	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
	"go.temporal.io/sdk/workflow"
)

func MyWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

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

	logger.Info("HelloWorld workflow completed")

	return "result", nil
}

func MyActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("MyActivity - reading /tmp/result")
	data, err := os.ReadFile("/tmp/result")
	return string(data), err
}
