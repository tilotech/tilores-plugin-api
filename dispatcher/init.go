package dispatcher

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
)

// KillFunc is used to shutdown the plugin once it is no longer used
type KillFunc func()

// Initialize will start the given plugin and return its Dispatcher
// implementation.
//
// The returned KillFunc must be executed once the plugin will no longer be
// used. This also applies in case of an error!
func Initialize(cmd *exec.Cmd) (Dispatcher, KillFunc, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"dispatcher": &Plugin{},
		},
		Cmd:        cmd,
		SyncStdout: os.Stdout,
		SyncStderr: os.Stderr,
	})

	rpcClient, err := client.Client()
	if err != nil {
		return nil, client.Kill, err
	}

	raw, err := rpcClient.Dispense("dispatcher")
	if err != nil {
		return nil, client.Kill, err
	}

	impl, ok := raw.(Dispatcher)
	if !ok {
		return nil, client.Kill, fmt.Errorf("received invalid client: %T", raw)
	}

	return impl, client.Kill, nil
}
