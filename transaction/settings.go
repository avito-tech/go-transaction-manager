package transaction

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

import "time"

// Settings is the configuration of the Manager.
// Preferable to implement as an immutable struct.
//
// settings.Settings is a default implementation of Settings.
type Settings interface {
	// EnrichBy fills nil properties from external Settings.
	EnrichBy(external Settings) Settings

	// CtxKey returns transaction.CtxKey for the transaction.Transaction.
	CtxKey() CtxKey
	CtxKeyOrNil() *CtxKey
	SetCtxKey(*CtxKey) Settings

	// Propagation returns transaction.Propagation.
	Propagation() Propagation
	PropagationOrNil() *Propagation
	SetPropagation(*Propagation) Settings

	// Cancelable defines that parent transaction.Transaction can cancel child transaction.Transaction or goroutines.
	Cancelable() bool
	CancelableOrNil() *bool
	SetCancelable(*bool) Settings

	// TimeoutOrNil returns time.Duration of the transaction.Transaction.
	TimeoutOrNil() *time.Duration
	SetTimeout(*time.Duration) Settings
}
