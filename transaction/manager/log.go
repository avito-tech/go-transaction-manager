package manager

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

type logger interface {
	Printf(format string, a ...interface{})
}

//nolint:gochecknoglobals // initializing default log, which does nothing
var defaultLog = log{}

type log struct{}

func (l log) Printf(_ string, _ ...interface{}) {}
