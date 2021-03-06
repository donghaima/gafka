package command

import (
	"testing"

	"github.com/funkygao/assert"
)

func TestSortBackendByName(t *testing.T) {
	all := []Backend{
		Backend{Name: "p2"},
		Backend{Name: "p1"},
		Backend{Name: "p3"},
	}

	r := sortBackendByName(all)
	t.Logf("%+v", r)
	assert.Equal(t, "p1", r[0].Name)
	assert.Equal(t, "p2", r[1].Name)
	assert.Equal(t, "p3", r[2].Name)
}
