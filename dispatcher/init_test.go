package dispatcher

import (
	"context"
	"os/exec"
)

func ExampleInitialize() {
	dispatcher, kill, err := Initialize(exec.Command("/path/to/plugin"))
	defer kill()
	if err != nil {
		panic(err)
	}
	dispatcher.Entity(context.Background(), &EntityInput{ID: "uuid"})
}
