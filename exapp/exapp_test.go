package exapp_test

import (
	"testing"

	"github.com/mteodor/edgex-app/exapp/events"

	"github.com/mteodor/edgex-app/exapp"
)

func TestExapp(t *testing.T) {
	cases := map[string]struct {
		err error
	}{
		"validate edgex-app": {events.ErrNotFound},
	}

	for desc, tc := range cases {
		err := exapp.Validate()
		//assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}
