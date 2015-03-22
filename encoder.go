package drum

import (
	"bytes"
	"os"
)

func EncodeToFile(path string, p *Pattern) error {

	fo, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fo.Close()

	bytes, err := EncodeToBytes(p)
	if err != nil {
		return err
	}
	if _, err := fo.Write(bytes); err != nil {
		return err
	}

	return nil
}

func EncodeToBytes(p *Pattern) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Initialize empty sliceFormat
	sf := &sliceFormat{}

	// Encode the slice header.
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
