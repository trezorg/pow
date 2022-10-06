package pow

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPow(t *testing.T) {
	req := Request{
		RemoteIP:   "127.0.0.1",
		RemotePort: 8080,
		SeedID:     123,
	}
	puzzle := Make(req)
	solution := Find(puzzle, 8)
	result := Check(req, solution, 8)
	require.True(t, result)
}

func TestMarshalUnMarshalPuzzle(t *testing.T) {
	puzzle := Puzzle{
		Solution: 1,
		SeedID:   1,
		Hash:     *(*[20]byte)([]byte("11111111111111111111")),
	}

	res, err := puzzle.MarshalBinary()
	require.NoError(t, err)
	newPuzzle := Puzzle{}
	err = newPuzzle.UnmarshalBinary(res)
	require.NoError(t, err)
	require.Equal(t, puzzle.Hash, newPuzzle.Hash)
	require.Equal(t, puzzle.Solution, newPuzzle.Solution)
	require.Equal(t, puzzle.SeedID, newPuzzle.SeedID)
}

func TestZeroes(t *testing.T) {
	data := []byte{255}
	require.False(t, zeroBits(data, 3))
	data = []byte{0}
	require.True(t, zeroBits(data, 8))
	data = []byte{0, 0}
	require.True(t, zeroBits(data, 16))
	data = []byte{4}
	require.False(t, zeroBits(data, 8))
	data = []byte{4}
	require.True(t, zeroBits(data, 5))
}

func TestReadPuzzleFromReader(t *testing.T) {
	puzzle := Puzzle{
		Solution: 1,
		SeedID:   1,
		Hash:     *(*[20]byte)([]byte("11111111111111111111")),
	}

	res, err := puzzle.MarshalBinary()
	require.NoError(t, err)
	newPuzzle := Puzzle{}
	r := bytes.NewReader(res)
	err = newPuzzle.Read(r)
	require.NoError(t, err)
	require.Equal(t, puzzle.Hash, newPuzzle.Hash)
	require.Equal(t, puzzle.Solution, newPuzzle.Solution)
	require.Equal(t, puzzle.SeedID, newPuzzle.SeedID)
}

func TestReadInvalidPuzzle(t *testing.T) {
	newPuzzle := Puzzle{}
	r := bytes.NewReader([]byte("test"))
	err := newPuzzle.Read(r)
	require.Error(t, err)
}
