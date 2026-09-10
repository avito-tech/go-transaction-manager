package main

import (
	"log"
)

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
	log.Printf("tr: %+v, isExp: %v", tr, ok)

	expTr := retExpTransaction()

	_, ok = expTr.(expTransaction)
	log.Printf("expTr: %+v, isExp: %v", expTr, ok)
}

func retTransaction() interface{} {
	return tr{}
}

func retExpTransaction() interface{} {
	return exTr{}
}
