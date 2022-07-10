package chinabond

import (
	"testing"
)

func TestYcDetail(t *testing.T) {
	d, err := YcDetail(6, "AAA", "2022-07-08")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(d)
}
