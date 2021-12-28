package dispatcher

import (
	"context"
	"fmt"

	api "github.com/tilotech/tilores-plugin-api"
	"github.com/tilotech/tilores-plugin-api/commons"
)

type server struct {
	impl Dispatcher
}

func (s *server) Entity(args map[string]interface{}, resp *api.Entity) error {
	ctx := context.Background() // TODO: replace with actual context
	id, err := commons.StringValue(args, "id")
	if err != nil {
		return err
	}
	entity, err := s.impl.Entity(ctx, id)
	if err != nil {
		return err
	}
	*resp = *entity
	return nil
}

func (s *server) Submit(args map[string]interface{}, resp *SubmitOutput) error {
	ctx := context.Background() // TODO: replace with actual context
	val, err := commons.Value(args, "input")
	if err != nil {
		return err
	}
	input, ok := val.(*SubmitInput)
	if !ok {
		return fmt.Errorf("key records is not a *SubmitInput but a %T", val)
	}
	submitOutput, err := s.impl.Submit(ctx, input)
	if err != nil {
		return err
	}
	*resp = *submitOutput
	return nil
}

func (s *server) Disassemble(args map[string]interface{}, resp *DisassembleOutput) error {
	ctx := context.Background() // TODO: replace with actual context
	val, err := commons.Value(args, "input")
	if err != nil {
		return err
	}
	input, ok := val.(*DisassembleInput)
	if !ok {
		return fmt.Errorf("key input is not a *DisassembleInput but a %T", val)
	}
	disassembleOutput, err := s.impl.Disassemble(ctx, input)
	if err != nil {
		return err
	}
	*resp = *disassembleOutput
	return nil
}

func (s *server) RemoveConnectionBan(args map[string]interface{}, _ *interface{}) error {
	ctx := context.Background() // TODO: replace with actual context
	val, err := commons.Value(args, "input")
	if err != nil {
		return err
	}
	input, ok := val.(*RemoveConnectionBanInput)
	if !ok {
		return fmt.Errorf("key input is not a *RemoveConnectionBanInput but a %T", val)
	}
	err = s.impl.RemoveConnectionBan(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) Search(args map[string]interface{}, resp *[]*api.Entity) error {
	ctx := context.Background() // TODO: replace with actual context
	val, err := commons.Value(args, "parameters")
	if err != nil {
		return err
	}
	parameters, ok := val.(*api.SearchParameters)
	if !ok {
		return fmt.Errorf("key parameters is not a *api.SearchParameters but a %T", val)
	}
	searchResult, err := s.impl.Search(ctx, parameters)
	if err != nil {
		return err
	}
	*resp = searchResult
	return nil
}
