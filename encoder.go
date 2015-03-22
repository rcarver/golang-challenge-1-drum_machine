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
	err := encodeSliceFormat(sf, p)
	if err != nil {
		return buf.Bytes(), err
	}

	// Accumulate tracks and their size.
	trackBytes := int64(0)
	tracks := make([]*trackFormat, 0, len(p.Tracks))

	// Encode each track.
	for _, t := range p.Tracks {
		tf := &trackFormat{}
		err := encodeTrackFormat(tf, t)
		if err != nil {
			return buf.Bytes(), err
		}
		trackBytes += tf.ByteSize()
		tracks = append(tracks, tf)
	}

	// Set the slice header's file size.
	sf.SetFileSize(trackBytes)

	// Write the sliceFormat.
	if err := sf.Encode(buf); err != nil {
		return buf.Bytes(), err
	}

	// Write all trackFormat.
	for _, tf := range tracks {
		if err := tf.Encode(buf); err != nil {
			return buf.Bytes(), err
		}
	}

	return buf.Bytes(), nil
}

func encodeSliceFormat(sf *sliceFormat, p *Pattern) error {
	sf.Magic = [13]byte{'S', 'P', 'L', 'I', 'C', 'E'}

	for i, c := range p.Version {
		sf.VersionBytes[i] = byte(c)
	}

	sf.Tempo = p.Tempo

	return nil
}

func encodeTrackFormat(tf *trackFormat, t *Track) error {
	tf.ID = uint32(t.ID)
	tf.NameSize = byte(len(t.Name))
	tf.Name = t.Name
	for i, s := range t.Steps {
		if s {
			tf.steps[i] = 1
		} else {
			tf.steps[i] = 0
		}

	}
	return nil
}
