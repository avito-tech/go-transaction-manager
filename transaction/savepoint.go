package transaction

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

import "fmt"

// This file copies an interface from https://github.com/DATA-DOG/go-txdb

// SavePoint defines the syntax to create savepoints
// within transaction.
type SavePoint interface {
	Create(id string) string
	Release(id string) string
	Rollback(id string) string
}

//revive:disable:exported
type DefaultSavePoint struct{}

// NewSavePoint is an implementation of transaction.SavePoint.
func NewSavePoint() DefaultSavePoint {
	return DefaultSavePoint{}
}

func (dsp DefaultSavePoint) Create(id string) string {
	return fmt.Sprintf("SAVEPOINT %s", id)
}

func (dsp DefaultSavePoint) Release(id string) string {
	return fmt.Sprintf("RELEASE SAVEPOINT %s", id)
}

func (dsp DefaultSavePoint) Rollback(id string) string {
	return fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", id)
}

//revive:enabled:exported
