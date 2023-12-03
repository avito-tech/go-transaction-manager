package trm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test")

func TestIsSkippable(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}

	tests := map[string]struct {
		args args
		want bool
	}{
		"skippable": {
			args: args{
				err: Skippable(Skippable(errTest)),
			},
			want: true,
		},
		"unSkippable": {
			args: args{
				err: errTest,
			},
			want: false,
		},
		"nil": {
			args: args{
				err: nil,
			},
			want: false,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := IsSkippable(tt.args.err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSkippable(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}

	tests := map[string]struct {
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		"skippable": {
			args: args{
				err: Skippable(Skippable(errTest)),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrSkip) &&
					assert.ErrorIs(t, err, errTest)
			},
		},
		"nil": {
			args: args{
				err: Skippable(Skippable(nil)),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := Skippable(tt.args.err)

			tt.wantErr(t, err)
		})
	}
}

func TestUnSkippable(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}

	tests := map[string]struct {
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		"unSkippable": {
			args: args{
				err: UnSkippable(UnSkippable(
					Skippable(Skippable(errTest)))),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotErrorIs(t, err, ErrSkip) &&
					assert.ErrorIs(t, err, errTest)
			},
		},
		"nil": {
			args: args{
				err: UnSkippable(UnSkippable(nil)),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := UnSkippable(tt.args.err)

			tt.wantErr(t, err)
		})
	}
}
