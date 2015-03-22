package drum

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestEncodeToBytes(t *testing.T) {
	// fixtures/pattern_1.splice
	p1 := &Pattern{}
	p1.Version = "0.808-alpha"
	p1.Tempo = 120.0
	p1.AddTrack(&Track{
		ID:   0,
		Name: "kick",
		Steps: [16]bool{
			true, false, false, false,
			true, false, false, false,
			true, false, false, false,
			true, false, false, false,
		},
	})
	p1.AddTrack(&Track{
		ID:   1,
		Name: "snare",
		Steps: [16]bool{
			false, false, false, false,
			true, false, false, false,
			false, false, false, false,
			true, false, false, false,
		},
	})
	p1.AddTrack(&Track{
		ID:   2,
		Name: "clap",
		Steps: [16]bool{
			false, false, false, false,
			true, false, true, false,
			false, false, false, false,
			false, false, false, false,
		},
	})
	p1.AddTrack(&Track{
		ID:   3,
		Name: "hh-open",
		Steps: [16]bool{
			false, false, true, false,
			false, false, true, false,
			true, false, true, false,
			false, false, true, false,
		},
	})
	p1.AddTrack(&Track{
		ID:   4,
		Name: "hh-close",
		Steps: [16]bool{
			true, false, false, false,
			true, false, false, false,
			false, false, false, false,
			true, false, false, true,
		},
	})
	p1.AddTrack(&Track{
		ID:   5,
		Name: "cowbell",
		Steps: [16]bool{
			false, false, false, false,
			false, false, false, false,
			false, false, true, false,
			false, false, false, false,
		},
	})
	tests := []struct {
		path    string
		pattern *Pattern
	}{
		{"pattern_1.splice", p1},
	}
	for _, test := range tests {
		encoded, err := EncodeToBytes(test.pattern)
		if err != nil {
			t.Fatalf("error encoding: %s", err)
		}
		fi, err := os.Open(path.Join("fixtures", test.path))
		if err != nil {
			t.Fatalf("error opening file %s: %s", test.path, err)
		}
		expEncoded, err := ioutil.ReadAll(fi)
		if err != nil {
			t.Fatalf("error reading file %s: %s", test.path, err)
		}
		if !bytes.Equal(encoded, expEncoded) {
			t.Fatalf("%s wasn't encoded as expected.\nGot:\n%s\nExpected:\n%s",
				test.path, hex.Dump(encoded), hex.Dump(expEncoded))
		}
	}

}
