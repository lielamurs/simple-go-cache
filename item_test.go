package sgcache

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestItemExpiry(t *testing.T) {
	tests := []struct {
		name        string
		givenItem   *Item
		wantOutcome bool
	}{
		{
			name: "Check expired item",
			givenItem: &Item{
				data: "Expired item",
				ttl:  time.Now().Add(-time.Minute),
			},
			wantOutcome: true,
		},
		{
			name: "Check active item",
			givenItem: &Item{
				data: "Active item",
				ttl:  time.Now().Add(time.Minute),
			},
			wantOutcome: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			gotOutcome := test.givenItem.expired()

			if diff := cmp.Diff(test.wantOutcome, gotOutcome); diff != "" {
				t.Errorf("outcome mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
