package ruffe

import "testing"

func TestMiddlewareWrapSequence(t *testing.T) {
	results := make([]bool, 3)
	err := NewMiddleware(HandlerFunc(func(_ Context) error {
		results[0] = true
		return nil
	})).Wrap(HandlerFunc(func(_ Context) error {
		if !results[0] {
			t.Fatalf("wrong sequence. first middleware wasn't triggered")
		}
		results[1] = true
		return nil
	})).Wrap(HandlerFunc(func(_ Context) error {
		if !results[1] {
			t.Fatalf("wrong sequence. second middleware wasn't triggered")
		}
		//results[2] = true
		return nil
	})).Handle(nil)
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
