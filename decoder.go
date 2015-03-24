package drum

import (
	"bufio"
	"io"
	"os"
)

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	// Open the file.
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	// Decode the data.
	return DecodePattern(fi)
}

// DecodePattern decodes the drum machine data accessed via reader.
func DecodePattern(reader io.Reader) (*Pattern, error) {
	buf := bufio.NewReader(reader)

	// Parse the overall slice into a Pattern.
	sf := &sliceFormat{}
	p, err := sf.DecodePattern(buf)
	if err != nil {
		return p, err
	}

	// Read the rest of the file for tracks.
	trackReader := io.LimitReader(buf, sf.TrackBytes()).(*io.LimitedReader)

	// Parse the tracks.
	for trackReader.N > 0 {
		tf := &trackFormat{}
		t, err := tf.DecodeTrack(trackReader)
		if err != nil {
			return p, err
		}
		p.AddTrack(t)
	}

	return p, nil
}
