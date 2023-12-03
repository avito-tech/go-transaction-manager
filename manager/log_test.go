package manager

import (
	"context"
	"testing"
)

func Test_log_Printf(t *testing.T) {
	t.Parallel()

	log{}.Warning(context.TODO(), "")
}
