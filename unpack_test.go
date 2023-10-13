package wxapkg

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	if err := Unpack("./tests/__APP__.wxapkg", "./unpack_out", false, nil); err != nil {
		t.Fatal(err)
	}
}
