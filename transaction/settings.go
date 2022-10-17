package transaction

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

import "time"

// Settings is the configuration of the Manager.
// Preferable to implement as an immutable struct.
type Settings interface {
	// EnrichBy fills nil properties from external Settings.
	EnrichBy(external Settings) Settings

	// TODO
	CtxKey() CtxKey
	CtxKeyOrNil() *CtxKey
	SetCtxKey(*CtxKey) Settings

	// TODO
	IsReadOnly() bool
	IsReadOnlyOrNil() *bool
	SetIsReadOnly(*bool) Settings

	Propagation() Propagation
	PropagationOrNil() *Propagation
	SetPropagation(*Propagation) Settings

	// TODO
	Timeout() time.Duration
	TimeoutOrNil() *time.Duration
	SetTimeout(*time.Duration) Settings
}
