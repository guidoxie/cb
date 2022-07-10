package jisilu

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Error(err)
	}
	idus, err := client.industryList()
	if err != nil {
		t.Error(err)
	}
	top := idus.TopLevel("490101")
	if top != nil {
		for _, r := range idus.SubLevel("77") {
			fmt.Println(r)
		}
	}
}
