package textutils_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/gsmcwhirter/mlai-workshop/pkg/textutils"
)

func TestNgramIterator(t *testing.T) {
	tests := []struct {
		name   string
		iter   *textutils.NgramIterator
		filter func([]byte) bool
		want   [][][]byte
	}{
		{
			name: "quick test",
			iter: textutils.NewNgramIterator(bytes.NewBufferString("foo bar baz quux"), 1, 3),
			want: [][][]byte{
				{[]byte("foo")},
				{[]byte("foo"), []byte("bar")},
				{[]byte("bar")},
				{[]byte("foo"), []byte("bar"), []byte("baz")},
				{[]byte("bar"), []byte("baz")},
				{[]byte("baz")},
				{[]byte("bar"), []byte("baz"), []byte("quux")},
				{[]byte("baz"), []byte("quux")},
				{[]byte("quux")},
			},
		},
		{
			name: "not enough data",
			iter: textutils.NewNgramIterator(bytes.NewBufferString("foo bar"), 1, 3),
			want: [][][]byte{
				{[]byte("foo")},
				{[]byte("foo"), []byte("bar")},
				{[]byte("bar")},
			},
		},
		{
			name: "exact 3 test",
			iter: textutils.NewNgramIterator(bytes.NewBufferString("foo bar baz quux"), 3, 3),
			want: [][][]byte{
				{[]byte("foo"), []byte("bar"), []byte("baz")},
				{[]byte("bar"), []byte("baz"), []byte("quux")},
			},
		},
		{
			name: "2 and 3 test",
			iter: textutils.NewNgramIterator(bytes.NewBufferString("foo bar baz quux"), 2, 3),
			want: [][][]byte{
				{[]byte("foo"), []byte("bar")},
				{[]byte("foo"), []byte("bar"), []byte("baz")},
				{[]byte("bar"), []byte("baz")},
				{[]byte("bar"), []byte("baz"), []byte("quux")},
				{[]byte("baz"), []byte("quux")},
			},
		},
		{
			name: "rewriting mins",
			iter: textutils.NewNgramIterator(bytes.NewBufferString("foo bar baz quux"), 0, 3),
			want: [][][]byte{
				{[]byte("foo")},
				{[]byte("foo"), []byte("bar")},
				{[]byte("bar")},
				{[]byte("foo"), []byte("bar"), []byte("baz")},
				{[]byte("bar"), []byte("baz")},
				{[]byte("baz")},
				{[]byte("bar"), []byte("baz"), []byte("quux")},
				{[]byte("baz"), []byte("quux")},
				{[]byte("quux")},
			},
		},
		{
			name: "rewriting maxs",
			iter: textutils.NewNgramIterator(bytes.NewBufferString("foo bar baz quux"), 3, 2),
			want: [][][]byte{
				{[]byte("foo"), []byte("bar"), []byte("baz")},
				{[]byte("bar"), []byte("baz"), []byte("quux")},
			},
		},
		{
			name: "filter starts with b",
			iter: textutils.NewNgramIterator(bytes.NewBufferString("foo bar foobar baz quux"), 1, 3),
			filter: func(b []byte) bool {
				if len(b) < 1 {
					return true
				}
				return b[0] != 'b'
			},
			want: [][][]byte{
				{[]byte("foo")},
				{[]byte("foo"), []byte("foobar")},
				{[]byte("foobar")},
				{[]byte("foo"), []byte("foobar"), []byte("quux")},
				{[]byte("foobar"), []byte("quux")},
				{[]byte("quux")},
			},
		},
		{
			name: "not enough data exact 3",
			iter: textutils.NewNgramIterator(bytes.NewBufferString("foo bar"), 3, 3),
			want: [][][]byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.iter.SetFilter(tt.filter)

			ct := 0

			for {
				gram, err := tt.iter.Next()
				if err != nil {
					if ct != len(tt.want) {
						t.Errorf("Not the right number of outputs; want = %v, got = %v", len(tt.want), ct)
					}
					break
				}

				if ct >= len(tt.want) {
					t.Errorf("NgramIterator.Next() got extra ngram %v", gram)
				} else if !reflect.DeepEqual(gram, tt.want[ct]) {
					t.Errorf("NgramIterator.Next() = %v, want %v", gram, tt.want[ct])
				}

				ct++
			}
		})
	}
}

var outgram [][]byte

func BenchmarkNgramIterator3(b *testing.B) {
	dat, err := ioutil.ReadFile("./testdata/aeneid.mb.txt")
	if err != nil {
		panic(err)
	}

	data := bytes.NewReader(dat)
	for n := 0; n < b.N; n++ {
		_, err := data.Seek(0, io.SeekStart)
		if err != nil {
			panic(err)
		}

		iter := textutils.NewNgramIterator(data, 1, 3)
		for {
			out, err := iter.Next()
			if err != nil {
				break
			}

			outgram = out
		}
	}
}

func BenchmarkNgramIterator4(b *testing.B) {
	dat, err := ioutil.ReadFile("./testdata/aeneid.mb.txt")
	if err != nil {
		panic(err)
	}

	data := bytes.NewReader(dat)
	for n := 0; n < b.N; n++ {
		_, err := data.Seek(0, io.SeekStart)
		if err != nil {
			panic(err)
		}

		iter := textutils.NewNgramIterator(data, 1, 4)
		for {
			out, err := iter.Next()
			if err != nil {
				break
			}

			outgram = out
		}
	}
}

func BenchmarkNgramIterator5(b *testing.B) {
	dat, err := ioutil.ReadFile("./testdata/aeneid.mb.txt")
	if err != nil {
		panic(err)
	}

	data := bytes.NewReader(dat)
	for n := 0; n < b.N; n++ {
		_, err := data.Seek(0, io.SeekStart)
		if err != nil {
			panic(err)
		}

		iter := textutils.NewNgramIterator(data, 1, 5)
		for {
			out, err := iter.Next()
			if err != nil {
				break
			}

			outgram = out
		}
	}
}
