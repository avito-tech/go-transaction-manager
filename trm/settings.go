package trm

//go:generate mockgen -source=$GOFILE -destination=drivers/mock/$GOFILE -package=mock

import (
	"time"
)

// Settings is the configuration of the Manager.
// Preferable to implement as an immutable struct.
//
// settings.Settings is a default implementation of Settings.
//
//nolint:interfacebloat
type Settings interface {
	// EnrichBy fills nil properties from external Settings.
	EnrichBy(external Settings) Settings

	// CtxKey returns trm.CtxKey for the trm.Transaction.
	CtxKey() CtxKey
	CtxKeyOrNil() *CtxKey
	SetCtxKey(*CtxKey) Settings

	// Propagation returns trm.Propagation.
	Propagation() Propagation
	PropagationOrNil() *Propagation
	SetPropagation(*Propagation) Settings

	// Cancelable defines that parent trm.Transaction can cancel child trm.Transaction or goroutines.
	Cancelable() bool
	CancelableOrNil() *bool
	SetCancelable(*bool) Settings

	// TimeoutOrNil returns time.Duration of the trm.Transaction.
	TimeoutOrNil() *time.Duration
	SetTimeout(*time.Duration) Settings
}
