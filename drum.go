// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import (
	"bytes"
	"fmt"
)

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	Version string
	Tempo   float32
	Tracks  []*Track
}

// AddTrack gives the pattern another track to play.
func (p *Pattern) AddTrack(t *Track) {
	p.Tracks = append(p.Tracks, t)
}

// String outputs a human readable view of the entire pattern.
func (p Pattern) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "Saved with HW Version: %s\n", p.Version)
	fmt.Fprintf(&buf, "Tempo: %g\n", p.Tempo)

	for _, track := range p.Tracks {
		fmt.Fprintf(&buf, "(%d) %s\t", track.ID, track.Name)

		fmt.Fprintf(&buf, "|")
		for i, step := range track.Steps {
			if step {
				fmt.Fprintf(&buf, "x")
			} else {
				fmt.Fprintf(&buf, "-")
			}
			if (i+1)%4 == 0 {
				fmt.Fprintf(&buf, "|")
			}
		}
		fmt.Fprintf(&buf, "\n")
	}

	return buf.String()
}

// Track is the pattern played by a single instrument within a Pattern.
type Track struct {
	ID    int
	Name  string
	Steps [16]bool
}
