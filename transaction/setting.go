package transaction

// Settings is configuration for Manager.
// TODO probably needs to separate Transaction and Manager settings.
type Settings interface {
	CtxKey() CtxKey
	IsReadOnly() bool
	Propagation() Propagation
}
