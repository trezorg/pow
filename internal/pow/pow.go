package pow

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"io"
	"math/rand"
	"net"
	"time"
)

type Puzzle struct {
	Hash     [20]byte
	Solution uint64
	SeedID   uint64
}

func (puzzle *Puzzle) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, *puzzle); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (puzzle *Puzzle) UnMarshal(buf []byte) error {
	if err := binary.Read(bytes.NewBuffer(buf), binary.BigEndian, puzzle); err != nil {
		return err
	}
	return nil
}

func (puzzle *Puzzle) read(r io.Reader) ([]byte, error) {
	b := make([]byte, 36)
	_, err := io.ReadFull(r, b)
	return b, err
}

func (puzzle *Puzzle) Read(r io.Reader) error {
	data, err := puzzle.read(r)
	if err != nil {
		return err
	}
	return puzzle.UnMarshal(data)
}

type Request struct {
	RemoteIP   string
	RemotePort uint16
	SeedID     uint64
}

func (req Request) hash() []byte {
	hashes := sha1.New()
	// ip - 4 + port - 2 + seed 8
	bs := make([]byte, 14)
	ip := net.ParseIP(req.RemoteIP)
	copy(bs[:4], ip.To4())
	binary.BigEndian.PutUint16(bs[4:], req.RemotePort)
	binary.BigEndian.PutUint64(bs[6:], req.SeedID)
	hashes.Write(bs)
	return hashes.Sum(nil)
}

func zeroBits(bytes []byte, n int) bool {
	zeroes := 0
	for _, b := range bytes {
		for i := 0; i < 8; i++ {
			bit := (b >> uint(7-i)) & 0x01
			if bit == 0 {
				zeroes++
			}
			if bit == 1 || zeroes == n {
				return zeroes == n
			}
		}
	}
	return zeroes == n
}

func Make(req Request) Puzzle {
	puz := Puzzle{}
	puz.Hash = *(*[20]byte)(req.hash())
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
	return zeroBits(res, zeroes)
}

func Find(puzzle Puzzle, zeroes int) Puzzle {
	buf := make([]byte, len(puzzle.Hash)+8)
	copy(buf, puzzle.Hash[:])
	rand.Seed(time.Now().UnixNano())
	for {
		solution := rand.Uint64()
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, solution)
		copy(buf[len(puzzle.Hash):], bs)
		hashes := sha1.New()
		hashes.Write(buf)
		res := hashes.Sum(nil)
		if zeroBits(res, zeroes) {
			return Puzzle{
				Solution: solution,
				SeedID:   puzzle.SeedID,
				Hash:     puzzle.Hash,
			}
		}
	}
}
