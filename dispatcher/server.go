package dispatcher

import (
	"context"
	"fmt"
	"time"

	"github.com/tilotech/tilores-plugin-api/commons"
)

type server struct {
	impl Dispatcher
}

func (s *server) Entity(args map[string]interface{}, resp *EntityOutput) error {
	ctx, cancel := createContext(args)
	defer cancel()

	val, err := commons.Value(args, "input")
	if err != nil {
		return err
	}
	input, ok := val.(*EntityInput)
	if !ok {
		return fmt.Errorf("key records is not a *EntityInput but a %T", val)
	}
	output, err := s.impl.Entity(ctx, input)
	if err != nil {
		return err
	}
	*resp = *output
	return nil
}

func (s *server) Submit(args map[string]interface{}, resp *SubmitOutput) error {
	ctx, cancel := createContext(args)
	defer cancel()

	val, err := commons.Value(args, "input")
	if err != nil {
		return err
	}
	input, ok := val.(*SubmitInput)
	if !ok {
		return fmt.Errorf("key records is not a *SubmitInput but a %T", val)
	}
	output, err := s.impl.Submit(ctx, input)
	if err != nil {
		return err
	}
	*resp = *output
	return nil
}

func (s *server) Disassemble(args map[string]interface{}, resp *DisassembleOutput) error {
	ctx, cancel := createContext(args)
	defer cancel()

	val, err := commons.Value(args, "input")
	if err != nil {
		return err
	}
	input, ok := val.(*DisassembleInput)
	if !ok {
		return fmt.Errorf("key input is not a *DisassembleInput but a %T", val)
	}
	output, err := s.impl.Disassemble(ctx, input)
	if err != nil {
		return err
	}
	*resp = *output
	return nil
}

func (s *server) RemoveConnectionBan(args map[string]interface{}, _ *interface{}) error {
	ctx, cancel := createContext(args)
	defer cancel()

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

func (s *server) Search(args map[string]interface{}, resp *SearchOutput) error {
	ctx, cancel := createContext(args)
	defer cancel()

	val, err := commons.Value(args, "input")
	if err != nil {
		return err
	}
	input, ok := val.(*SearchInput)
	if !ok {
		return fmt.Errorf("key input is not a *SearchInput but a %T", val)
	}
	output, err := s.impl.Search(ctx, input)
	if err != nil {
		return err
	}
	*resp = *output
	return nil
}

func createContext(args map[string]interface{}) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	if deadlineVal, err := commons.Value(args, "deadline"); err == nil {
		if deadline, ok := deadlineVal.(time.Time); ok {
			ctx, cancel := context.WithDeadline(ctx, deadline)
			return ctx, cancel
		}
	}
	return ctx, func() {}
}
