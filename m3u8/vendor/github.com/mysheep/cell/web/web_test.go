package web

import (
	"bytes"
	"testing"
)

func TestCompare(t *testing.T) {
	xs := []struct {
		givenA []byte
		givenB []byte
		want   bool
	}{
		{[]byte{1, 2, 3, 4}, []byte{1, 2, 3, 4}, true},  // same content
		{[]byte{1, 2, 3, 4}, []byte{4, 3, 2, 1}, false}, // different content
		{[]byte{1, 2, 3, 4}, []byte{1, 2, 3}, false},    // different length
		{[]byte{}, []byte{}, true},                      // both empty
	}

	for _, x := range xs {
		got := (bytes.Compare(x.givenA, x.givenB) == 0)
		if got != x.want {
			t.Errorf("Compare want: %v, got: %v", x.want, got)
		}
	}
}

func TestDownloader(t *testing.T) {
	// given
	in := make(chan string)
	out := make(chan []byte)

	xs := []struct {
		given string
		want  []byte
	}{
		{"http://test.com/masterplaylist.m3u8", []byte{}},
		{"http://test.com/mediaplaylist.m3u8", []byte{}},
	}

	// TODO: Fill content
	/*
		content, err := ioutil.ReadFile("testdata/masterplaylist.m3u8")
		if err != nil {
			panic(err)
		}
	*/

	fakeG := FakeGetter{}
	go Downloader(fakeG, in, out)

	for _, x := range xs {
		// when
		in <- x.given
		got := <-out

		// then
		if bytes.Compare(got, x.want) != 0 {
			t.Errorf("Downloader was incorrect, got:%v, want:%v", got, x.want)
		}
	}

}
