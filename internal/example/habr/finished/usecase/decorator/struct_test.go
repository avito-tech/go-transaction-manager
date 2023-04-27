package decorator

import (
	"context"
	"fmt"
)

type usecase struct{}

func (u usecase) Handle(ctx context.Context, in string) (string, error) {
	fmt.Println("inside", in)

	return "result", nil
}

func ExampleFMTDecorator() {
	uc := FMTDecorator[string, string](usecase{})

	fmt.Println(uc.Handle(context.Background(), "in"))

	// Output: start
	//inside in
	//finish
	//result <nil>
}
