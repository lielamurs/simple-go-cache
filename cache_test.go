package sgcache

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCache(t *testing.T) {
	cache := New(time.Duration(time.Second))
	defer cache.Close()

	tests := []struct {
		name        string
		givenKey    string
		givenEntry  string
		givenTTL    time.Duration
		setEntry    bool
		wantEntry   interface{}
		wantOutcome bool
	}{
		{
			name:        "Test active entry",
			givenKey:    "active",
			givenEntry:  "1 minute entry",
			givenTTL:    time.Minute,
			setEntry:    true,
			wantEntry:   "1 minute entry",
			wantOutcome: true,
		},
		{
			name:        "Test unset entry",
			givenKey:    "unset",
			givenEntry:  "not set entry",
			givenTTL:    time.Minute,
			setEntry:    false,
			wantEntry:   nil,
			wantOutcome: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			if test.setEntry == true {
				cache.Set(test.givenKey, test.givenEntry, test.givenTTL)
			}

			gotEntry, gotOutcome := cache.Get(test.givenKey)

			if diff := cmp.Diff(test.wantOutcome, gotOutcome); diff != "" {
				t.Errorf("outcome mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(test.wantEntry, gotEntry); diff != "" {
				t.Errorf("outcome mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestCacheDelete(t *testing.T) {
	cache := New(time.Duration(time.Second))
	defer cache.Close()

	givenKey, givenEntry := "deletable", "Deletable entry"
	cache.Set(givenKey, givenEntry, time.Minute)

	gotEntry, gotOutcome := cache.Get(givenKey)
	if diff := cmp.Diff(true, gotOutcome); diff != "" {
		t.Errorf("outcome mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(givenEntry, gotEntry); diff != "" {
		t.Errorf("outcome mismatch (-want +got):\n%s", diff)
	}

	cache.Delete(givenKey)
	gotEntry, gotOutcome = cache.Get(givenKey)
	if diff := cmp.Diff(false, gotOutcome); diff != "" {
		t.Errorf("outcome mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(nil, gotEntry); diff != "" {
		t.Errorf("outcome mismatch (-want +got):\n%s", diff)
	}
}

func TestCleanup(t *testing.T) {
	cache := New(time.Duration(time.Millisecond))
	defer cache.Close()

	givenKey, givenEntry := "cleanup", "cleanup entry"
	cache.Set(givenKey, givenEntry, time.Second)

	gotEntry, gotOutcome := cache.Get(givenKey)
	if diff := cmp.Diff(true, gotOutcome); diff != "" {
		t.Errorf("outcome mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(givenEntry, gotEntry); diff != "" {
		t.Errorf("outcome mismatch (-want +got):\n%s", diff)
	}

	time.Sleep(time.Second * 2)
	gotEntry, gotOutcome = cache.Get(givenKey)
	if diff := cmp.Diff(false, gotOutcome); diff != "" {
		t.Errorf("outcome mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(nil, gotEntry); diff != "" {
		t.Errorf("outcome mismatch (-want +got):\n%s", diff)
	}
}
