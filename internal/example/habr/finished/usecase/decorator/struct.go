// detailed description  https://threedots.tech/post/increasing-cohesion-in-go-with-generic-decorators/
package decorator

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm"
)

type Usecase[In any, Out any] interface {
	Handle(ctx context.Context, in In) (Out, error)
}

type txDecorator[In any, Out any] struct {
	manager trm.Manager
	usecase Usecase[In, Out]
}

func TxDecorate[In any, Out any](m trm.Manager, u Usecase[In, Out]) Usecase[In, Out] {
	return &txDecorator[In, Out]{manager: m, usecase: u}
}

func (d *txDecorator[In, Out]) Handle(ctx context.Context, in In) (out Out, err error) {
	var emptyOut Out

	err = d.manager.Do(ctx, func(ctx context.Context) error {
		out, err = d.usecase.Handle(ctx, in)

		return err
	})
	if err != nil {
		return emptyOut, err
	}

	return out, nil
}

type fmtDecorator[In any, Out any] struct {
	usecase Usecase[In, Out]
}

func FMTDecorator[In any, Out any](u Usecase[In, Out]) Usecase[In, Out] {
	return &fmtDecorator[In, Out]{usecase: u}
}

func (d *fmtDecorator[In, Out]) Handle(ctx context.Context, in In) (out Out, err error) {
	fmt.Println("start")
	defer fmt.Println("finish")

	return d.usecase.Handle(ctx, in)
}
