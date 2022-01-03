package dispatcher

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Handshake is used to verify plugin compatibility
//
// This is a user experience feature, not a security feature!
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "DISPATCHER_PLUGIN",
	MagicCookieValue: "opsPsgZd3qsXpjgj69j5",
}

// Plugin includes Dispatcher implementation
type Plugin struct {
	Impl Dispatcher
}

// Server the server side of the plugin
func (p *Plugin) Server(_ *plugin.MuxBroker) (interface{}, error) {
	return &server{impl: p.Impl}, nil
}

// Client the client side of the plugin
func (*Plugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &client{client: c}, nil
}
