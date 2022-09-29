package pow

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPow(t *testing.T) {
	req := Request{
		RemoteIP:   "127.0.0.1",
		RemotePort: 8080,
		SeedID:     123,
	}
	puzzle := Make(req)
	solution := Find(puzzle, 3)
	result := Check(req, solution, 3)
	assert.True(t, result)
}
