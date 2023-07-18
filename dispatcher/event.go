package dispatcher

import (
	"encoding/json"
	"fmt"

	api "github.com/tilotech/tilores-plugin-api"
)

// AssembleEvent is used to parse a Kinesis Event that contains the data for
// the assemble lambda.
//
// The event can contain different types of payload depending on the action that
// needs to be performed.
//
// In case of type "ASSEMBLE", the payload contains []*api.Record entries.
//
// In case of type "DISASSEMBLE", the payload is a *dispatcher.DisassembleInput.
//
// For backwards compatibility, the old event input for assemble requests is
// also supported. The output after unmarshalling will be the same as for type
// "ASSEMBLE".
type AssembleEvent struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

const (
	// EventTypeAssemble is used when the payload is for the assemble process.
	EventTypeAssemble = "ASSEMBLE"

	// EventTypeDisassemble is used when the payload is for the disassemble process.
	EventTypeDisassemble = "DISASSEMBLE"
)

// UnmarshalJSON parses the provided bytes and populates the AssembleEvent.
func (r *AssembleEvent) UnmarshalJSON(b []byte) error {
	partial := &struct {
		Type    string
		Payload json.RawMessage
	}{}
	err := json.Unmarshal(b, partial)
	if err != nil {
		partial.Type = EventTypeAssemble
		partial.Payload = b
	}
	var payload any
	switch partial.Type {
	case EventTypeAssemble:
		pl := []*api.Record{}
		err = json.Unmarshal(partial.Payload, &pl)
		payload = pl
	case EventTypeDisassemble:
		payload = &DisassembleInput{}
		err = json.Unmarshal(partial.Payload, payload)
	default:
		return fmt.Errorf("invalid type %s", partial.Type)
	}
	if err != nil {
		return err
	}
	r.Type = partial.Type
	r.Payload = payload
	return nil
}
