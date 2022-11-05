package manager

import "context"

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

type logger interface {
	Warning(ctx context.Context, msg string)
}

//nolint:gochecknoglobals // initializing default log, which does nothing
var defaultLog = log{}

type log struct{}

func (l log) Warning(context.Context, string) {}
