package ruffe

import (
	"net/http"
	"testing"
)

type fakeCtx struct {
	http.ResponseWriter
}

func (f fakeCtx) done() bool                    { return false }
func (f fakeCtx) Request() *http.Request        { return nil }
func (f fakeCtx) Bind(interface{}) error        { return nil }
func (f fakeCtx) Result(int, interface{}) error { return nil }

func TestMiddlewareWrapSequence(t *testing.T) {
	results := make([]bool, 3)
	err := NewMiddleware(HandlerFunc(func(_ Context) error {
		results[0] = true
		return nil
	})).BeforeFunc(func(_ Context) error {
		if !results[0] {
			t.Fatalf("wrong sequence. first middleware wasn't triggered")
		}
		results[1] = true
		return nil
	}).BeforeFunc(func(_ Context) error {
		if !results[1] {
			t.Fatalf("wrong sequence. second middleware wasn't triggered")
		}
		results[2] = true
		return nil
	}).Handle(fakeCtx{})
	if err != nil {
		t.FailNow()
	}
	for i, r := range results {
		if !r {
			t.Errorf("middleware %d wasn't triggered", i)
			break
		}
	}
}
