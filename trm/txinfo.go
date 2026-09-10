package trm

// TxInfo exposes transaction metadata for hooks.
// Implementations must be thread-safe; reads during nested scopes may observe transient state.
type TxInfo interface {
	Propagation() Propagation
	IsNew() bool
	IsNested() bool
	NestingLevel() int
}
