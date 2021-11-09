package dispatcher

import (
	"context"

	api "gitlab.com/tilotech/tilores-plugin-api"
	"gitlab.com/tilotech/tilores-plugin-api/commons"
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
