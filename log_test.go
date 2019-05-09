package NetWorkRate

import (
	"testing"
)

func TestExampleLogger(t *testing.T) {
	r, _ := FastGet(false, nil, 1)

	l := NewCustomLogger("test.log", true, true, 15)
	l.Println(r)
}
