package manager

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

type logger interface {
	Printf(format string, a ...interface{})
}

// WithLog sets logger for Manager.
func WithLog(l logger) Opt {
	return func(m *Manager) {
		m.log = l
	}
}

//nolint:gochecknoglobals // initializing default log, which does nothing
var defaultLog = log{}

type log struct{}

func (l log) Printf(_ string, _ ...interface{}) {}
