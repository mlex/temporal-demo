# Demo for SRE Meetup Munich in July 2024

See [slides](Rollout-Automation-for-ClickHouse-Cloud-with-Temporal.pdf)

# Steps to run this sample:
1) Run a [Temporal service](https://github.com/temporalio/samples-go/tree/main/#how-to-use).
2) Run the following command to start the worker
```
go run worker/main.go
```
3) Run the following command to start the example
```
go run starter/main.go
```
