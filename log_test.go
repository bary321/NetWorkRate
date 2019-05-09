package NetWorkRate

import (
	"testing"
)

func TestExampleLogger(t *testing.T) {
	r, _ := FastGet(false, nil, 1)

	l := NewCustomLogger("test.log")
	l.Println(r)
}
