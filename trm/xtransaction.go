package trm

// XTransaction is an extended transaction that wraps trm.Transaction and adds TxInfo.
// It is stored in context.Context instead of the raw transaction when using XDo/XDoWithSettings.
type XTransaction interface {
	Transaction
	TxInfo() TxInfo
}
