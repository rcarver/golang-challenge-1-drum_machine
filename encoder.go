package drum

import (
	"bytes"
	"io/ioutil"
)

// EncodeToFile takes Pattern p and writes it to path.
func EncodeToFile(path string, p *Pattern) error {
	bytes, err := EncodeToBytes(p)
	if err != nil {
		return err
	}
	// TODO: use os.FileMode?
	if err := ioutil.WriteFile(path, bytes, 0x755); err != nil {
		return err
	}
	return nil
}

// EncodeToBytes takes Pattern p and returns its binary representation.
func EncodeToBytes(p *Pattern) ([]byte, error) {
	buf := &bytes.Buffer{}

	// Encode the slice header.
	sf := &sliceFormat{}
	err := sf.EncodePattern(p)
	if err != nil {
		return buf.Bytes(), err
	}

	// Accumulate tracks and their size.
	trackBytes := int64(0)
	tracks := make([]*trackFormat, 0, len(p.Tracks))

	// Encode each track.
	for _, t := range p.Tracks {
		tf := &trackFormat{}
		err := tf.EncodeTrack(t)
		if err != nil {
			return buf.Bytes(), err
		}
		trackBytes += tf.ByteSize()
		tracks = append(tracks, tf)
	}

	// Set the slice header's file size.
	sf.SetFileSize(trackBytes)

	// Write the sliceFormat.
	if err := sf.Write(buf); err != nil {
		return buf.Bytes(), err
	}

	// Write all trackFormat.
	for _, tf := range tracks {
		if err := tf.Write(buf); err != nil {
			return buf.Bytes(), err
		}
	}

	return buf.Bytes(), nil
}
