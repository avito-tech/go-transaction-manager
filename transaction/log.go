package transaction

type logger interface {
	Printf(format string, a ...interface{})
}

// WithLog sets logger for ManagerImpl.
func WithLog(l logger) ManagerOpt {
	return func(m *ManagerImpl) {
		if l == nil {
			l = defaultLog
		}

		m.log = l
	}
}

//nolint:gochecknoglobals // initializing default log, which does nothing
var defaultLog = log{}

type log struct{}

func (l log) Printf(_ string, _ ...interface{}) {}
