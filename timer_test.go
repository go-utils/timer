package timer

import (
	"errors"
	"testing"
	"time"
)

// testContext - test context
type testContext struct {
	Ch  chan struct{}
	err error
}

// newTestContext - constructor
func newTestContext() *testContext {
	return &testContext{
		Ch: make(chan struct{}),
	}
}

// Deadline ...
func (tc *testContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

// Done ...
func (tc *testContext) Done() <-chan struct{} {
	return tc.Ch
}

// Err ...
func (tc *testContext) Err() error {
	return tc.err
}

// Value ...
func (tc *testContext) Value(_ interface{}) interface{} {
	return nil
}

// SetError ...
func (tc *testContext) SetError(err error) {
	tc.err = err
}

func TestTimer(t *testing.T) {
	type args struct {
		c *testContext
		d time.Duration
		e error
		f func(ctx *testContext)
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				c: newTestContext(),
				d: 100 * time.Millisecond,
				e: nil,
				f: nil,
			},
			wantErr: false,
		},
		{
			name: "success_done",
			args: args{
				c: newTestContext(),
				d: 100 * time.Second,
				e: errors.New("test error"),
				f: func(ctx *testContext) {
					close(ctx.Ch)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // escape: Using the variable on range scope `tt` in loop literal.
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.e != nil {
				tt.args.c.SetError(tt.args.e)
			}

			if tt.args.f != nil {
				tt.args.f(tt.args.c)
				if err := RunWithContext(tt.args.c, tt.args.d); (err != nil) != tt.wantErr {
					t.Errorf("RunWithContext() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err := Run(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
