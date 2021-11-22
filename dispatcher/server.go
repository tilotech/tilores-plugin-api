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

func (s *server) Submit(args map[string]interface{}, resp *SubmissionResult) error {
	ctx := context.Background() // TODO: replace with actual context
	val, err := commons.Value(args, "records")
	if err != nil {
		return err
	}
	records, ok := val.([]*api.Record)
	if !ok {
		return fmt.Errorf("key records is not a records list but a %T", val)
	}
	submissionResult, err := s.impl.Submit(ctx, records)
	if err != nil {
		return err
	}
	*resp = *submissionResult
	return nil
}

func (s *server) Search(args map[string]interface{}, resp *[]*api.Entity) error {
	ctx := context.Background() // TODO: replace with actual context
	val, err := commons.Value(args, "parameters")
	if err != nil {
		return err
	}
	parameters, ok := val.(map[string]interface{})
	if !ok {
		return fmt.Errorf("key parameters is not a map[string]interface{} but a %T", val)
	}
	searchResult, err := s.impl.Search(ctx, parameters)
	if err != nil {
		return err
	}
	*resp = searchResult
	return nil
}
