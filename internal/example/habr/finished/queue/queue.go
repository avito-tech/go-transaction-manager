package queue

import (
	"context"
)

type Queue[In any] struct{}

func (Queue[In]) Publish(_ context.Context, in In) error {
	return nil
}
