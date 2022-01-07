package dispatcher_test

import (
	"context"
	"os/exec"

	"github.com/tilotech/tilores-plugin-api/dispatcher"
)

func ExampleInitialize() {
	dsp, kill, rc, err := dispatcher.Initialize(exec.Command("/path/to/plugin"), nil)
	defer kill()
	if err != nil {
		panic(err)
	}
	dsp.Entity(context.Background(), &dispatcher.EntityInput{ID: "uuid"})

	dsp, kill, _, err = dispatcher.Initialize(exec.Command("/path/to/plugin"), rc)
	defer kill()
	if err != nil {
		panic(err)
	}

	dsp.Entity(context.Background(), &dispatcher.EntityInput{ID: "another-uuid"})
}
