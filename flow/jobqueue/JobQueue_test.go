package jobqueue_test

import (
	"fmt"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/flow/jobqueue"
	"testing"
)

func TestJobQueue(t *testing.T) {
	work := func(input interface{}) interface{} {
		text := input.(string)
		return len(text)
	}

	jobs := jobqueue.New(work)

	jobs.Queue("Hello World")
	jobs.Queue("Test")

	results := jobs.Wait()

	fmt.Println(conv.String(results))
}
