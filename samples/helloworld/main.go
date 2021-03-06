package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CodeLinkIO/goworkers/tasks"
	"github.com/CodeLinkIO/goworkers/triggers/notifiers"
	"github.com/CodeLinkIO/goworkers/workers"
)

func handleHello(ctx context.Context, task tasks.Task) error {
	workerContext, ok := ctx.Value(workers.WorkerContextKey).(workers.WorkerContext)
	if !ok {
		return fmt.Errorf("Cannot get worker context")
	}
	processingTask, ok := task.(*helloTask)
	if !ok {
		return fmt.Errorf("Cannot parse task")
	}

	fmt.Println(fmt.Sprintf("%v: Handle task %v %v", workerContext.ID, processingTask.Type(), processingTask.Message))
	time.Sleep(200 * time.Millisecond)
	return nil
}

func main() {
	options := workers.ControllerOptions{
		NumOfWorker: 10,
	}
	ctx := context.Background()

	router := workers.NewRouter()
	router.Register("Hello", handleHello)

	trigger := notifiers.NewTrigger()

	controller := workers.NewController(ctx, trigger, router, options)

	go func() {
		for i := 0; i < 20; i++ {
			trigger.Notify(&helloTask{Message: fmt.Sprintf("tsk_%d", i)})
			time.Sleep(100 * time.Millisecond)
		}

		time.Sleep(time.Second)
		controller.Stop()
	}()

	err := controller.Run()
	if err != nil {
		log.Fatal(err)
	}
}
