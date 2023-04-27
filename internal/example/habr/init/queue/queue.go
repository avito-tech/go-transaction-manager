package queue

import "init/domain"

type Queue[In any] struct{}

func (Queue[In]) Publish(_ domain.Tr, in In) error {
	return nil
}
