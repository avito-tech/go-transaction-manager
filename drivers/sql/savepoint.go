package sql

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

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

// NewSavePoint is an implementation of trm.SavePoint.
func NewSavePoint() DefaultSavePoint {
	return DefaultSavePoint{}
}

func (dsp DefaultSavePoint) Create(id string) string {
	return "SAVEPOINT " + id
}

func (dsp DefaultSavePoint) Release(id string) string {
	return "RELEASE SAVEPOINT " + id
}

func (dsp DefaultSavePoint) Rollback(id string) string {
	return "ROLLBACK TO SAVEPOINT " + id
}

//revive:enabled:exported
