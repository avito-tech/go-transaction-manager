package main

import "fmt"

type transaction interface {
	isTransaction()
}

type expTransaction interface {
	transaction
	isExpTransaction()
}

type tr struct{}

func (t tr) isTransaction() {}

type exTr struct{}

func (e exTr) isTransaction() {}

func (e exTr) isExpTransaction() {}

func main() {
	tr := retTransaction()

	_, ok := tr.(expTransaction)
	fmt.Printf("tr: %+v, isExp: %v\n", tr, ok)

	expTr := retExpTransaction()

	_, ok = expTr.(expTransaction)
	fmt.Printf("expTr: %+v, isExp: %v\n", expTr, ok)
}

func retTransaction() transaction {
	return tr{}
}

func retExpTransaction() expTransaction {
	return exTr{}
}
