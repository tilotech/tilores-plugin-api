package dispatcher

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
)

// KillFunc is used to shutdown the plugin once it is no longer used
type KillFunc func()

// Initialize will start the given plugin and return its Dispatcher implementation
// or reattach to the already initialized plugin if reattachConfig is provided.
//
// The returned KillFunc must be executed once the plugin will no longer be
// used. This also applies in case of an error!
func Initialize(cmd *exec.Cmd, reattachConfig *plugin.ReattachConfig) (Dispatcher, KillFunc, *plugin.ReattachConfig, error) {
	if reattachConfig != nil {
		dsp, killFunc, rc, err := reattach(reattachConfig)
		if err == nil {
			return dsp, killFunc, rc, nil
		}
	}

	return start(cmd)
}

func start(cmd *exec.Cmd) (Dispatcher, KillFunc, *plugin.ReattachConfig, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"dispatcher": &Plugin{},
		},
		Cmd:        cmd,
		SyncStdout: os.Stdout,
		SyncStderr: os.Stderr,
	})

	return createDispatcher(client)
}

func reattach(reattachConfig *plugin.ReattachConfig) (Dispatcher, KillFunc, *plugin.ReattachConfig, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"dispatcher": &Plugin{},
		},
		Reattach:   reattachConfig,
		SyncStdout: os.Stdout,
		SyncStderr: os.Stderr,
	})

	return createDispatcher(client)
}

func createDispatcher(client *plugin.Client) (Dispatcher, KillFunc, *plugin.ReattachConfig, error) {
	rpcClient, err := client.Client()
	if err != nil {
		return nil, client.Kill, nil, err
	}

	raw, err := rpcClient.Dispense("dispatcher")
	if err != nil {
		return nil, client.Kill, nil, err
	}

	impl, ok := raw.(Dispatcher)
	if !ok {
		return nil, client.Kill, nil, fmt.Errorf("received invalid client: %T", raw)
	}

	return impl, client.Kill, client.ReattachConfig(), nil
}
