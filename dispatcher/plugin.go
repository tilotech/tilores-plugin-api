package dispatcher

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "DISPATCHER_PLUGIN",
	MagicCookieValue: "opsPsgZd3qsXpjgj69j5", // UX feature, not security feature!
}

type Plugin struct {
	Impl Dispatcher
}

func (p *Plugin) Server(_ *plugin.MuxBroker) (interface{}, error) {
	return &server{impl: p.Impl}, nil
}

func (*Plugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &client{client: c}, nil
}
