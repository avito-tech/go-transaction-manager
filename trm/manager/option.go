package manager

import "github.com/avito-tech/go-transaction-manager/trm/v2"

// WithLog sets logger for Manager.
func WithLog(l logger) Opt {
	return func(m *Manager) error {
		m.log = l

		return nil
	}
}

// WithSettings sets trm.Settings for Manager.
func WithSettings(s trm.Settings) Opt {
	return func(m *Manager) error {
		m.settings = s

		return nil
	}
}

// WithCtxManager sets trm.Settings for Manager.
func WithCtxManager(c trm.CtxManager) Opt {
	return func(m *Manager) error {
		m.ctxManager = c

		return nil
	}
}
