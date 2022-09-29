package pow

import (
	"crypto/sha1"
	"encoding/binary"
	"math/rand"
	"strconv"
	"time"
)

type Puzzle struct {
	Solution uint64
	SeedID   int
	Hash     []byte
}

type Request struct {
	RemoteIP   string
	RemotePort int
	SeedID     int
}

func (req Request) hash() []byte {
	var buf []byte
	hashes := sha1.New()
	buf = append(buf, []byte(req.RemoteIP)...)
	buf = append(buf, []byte(strconv.Itoa(req.RemotePort))...)
	buf = append(buf, []byte(strconv.Itoa(req.SeedID))...)
	hashes.Write(buf)
	return hashes.Sum(nil)
}

func isZeroFirst(bytes []byte, n int) bool {
	b := byte(0)
	for _, s := range bytes[:n] {
		b |= s
	}
	return b == 0
}

func Make(req Request) Puzzle {
	puz := Puzzle{}
	puz.Hash = req.hash()
	puz.SeedID = req.SeedID
	return puz
}

func Check(req Request, puzzle Puzzle, zeroes int) bool {
	hash := req.hash()
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, puzzle.Solution)
	hash = append(hash, bs...)
	hashes := sha1.New()
	hashes.Write(hash)
	res := hashes.Sum(nil)
	return isZeroFirst(res, zeroes)
}

func Find(puzzle Puzzle, zeroes int) Puzzle {
	buf := make([]byte, len(puzzle.Hash)+8)
	copy(buf, puzzle.Hash)
	rand.Seed(time.Now().UnixNano())
	for {
		solution := rand.Uint64()
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, solution)
		copy(buf[len(puzzle.Hash):], bs)
		hashes := sha1.New()
		hashes.Write(buf)
		res := hashes.Sum(nil)
		if isZeroFirst(res, zeroes) {
			return Puzzle{
				Solution: solution,
				SeedID:   puzzle.SeedID,
				Hash:     puzzle.Hash,
			}
		}
	}
}
