package textutils

import (
	"bufio"
	"io"
)

type NgramIterator struct {
	s    *bufio.Scanner
	minN int
	maxN int

	currGrams [][]byte
	currMaxN  int
	currN     int

	filter func([]byte) bool
}

func NewNgramIterator(r io.Reader, minN, maxN int) *NgramIterator {
	if minN <= 0 {
		minN = 1
	}

	if minN > maxN {
		maxN = minN
	}

	iter := &NgramIterator{
		s:    bufio.NewScanner(r),
		minN: minN,
		maxN: maxN,
	}

	iter.s.Split(bufio.ScanWords)

	return iter
}

func (iter *NgramIterator) SetFilter(f func([]byte) bool) {
	iter.filter = f
}

func (iter *NgramIterator) refill() bool {
	var n int
	switch {
	case iter.currGrams == nil:
		iter.currGrams = make([][]byte, iter.maxN)
		n = 0
		iter.currMaxN = 1
	case iter.currMaxN < iter.maxN:
		n = iter.currMaxN
		iter.currMaxN++
	default:
		n = copy(iter.currGrams, iter.currGrams[1:])
	}

	var gram []byte // = nil
	for gram == nil {
		if !iter.s.Scan() {
			// fmt.Println("empty")
			return false
		}

		gram = iter.s.Bytes()
		if len(gram) == 0 {
			gram = nil
			continue
		}

		if iter.filter != nil && !iter.filter(gram) {
			gram = nil
		}
	}

	iter.currGrams[n] = iter.s.Bytes()
	// fmt.Printf("currGrams: %+v\n", iter.currGrams)
	iter.currN = n + 1
	return true
}

func (iter *NgramIterator) Next() ([][]byte, error) {
	// initial conditions
	for iter.currMaxN < iter.minN {
		if !iter.refill() {
			return nil, io.EOF
		}
	}

	// actual next logic
	if iter.currN == 0 {
		// fmt.Println("refilling")
		if !iter.refill() {
			return nil, io.EOF
		}
	}

	// fmt.Printf("currMaxN=%v, currN=%v ", iter.currMaxN, iter.currN)
	out := iter.currGrams[iter.currMaxN-iter.currN : iter.currMaxN]
	// fmt.Printf("gram=%+v\n", out)

	iter.currN--
	if iter.currN < iter.minN {
		iter.currN = 0
	}

	return out, nil
}
