package manager

import "github.com/avito-tech/go-transaction-manager/transaction"

// WithLog sets logger for Manager.
func WithLog(l logger) Opt {
	return func(m *Manager) {
		m.log = l
	}
}

// WithSettings sets transaction.Settings for Manager.
func WithSettings(s transaction.Settings) Opt {
	return func(m *Manager) {
		m.settings = s
	}
}

// WithCtxManager sets transaction.Settings for Manager.
func WithCtxManager(с transaction.СtxManager) Opt {
	return func(m *Manager) {
		m.ctxManager = с
	}
}
