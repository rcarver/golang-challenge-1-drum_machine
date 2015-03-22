package drum

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
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
	p := &Pattern{}

	// Buffer all reading.
	buffer := bufio.NewReader(reader)

	// Parse the overall slice into a Pattern.
	sf := &sliceFormat{}
	err := sf.DecodePattern(p, buffer)
	if err != nil {
		return p, err
	}

	// Read the rest of the file for tracks.
	trackReader := io.LimitReader(buffer, sf.TrackBytes()).(*io.LimitedReader)

	// Parse the tracks.
	for trackReader.N > 0 {
		t := &Track{}
		tf := &trackFormat{}
		err := tf.DecodeTrack(t, trackReader)
		if err != nil {
			return p, err
		}
		p.AddTrack(t)
	}

	return p, nil
}

// sliceFormat is the low level binary format for the slice.
type sliceFormat struct {
	Magic        [13]byte
	FileSize     byte
	VersionBytes [32]byte
	Tempo        float32
}

// DecodePattern reads binary from reader and applies it to the Pattern.
func (sf *sliceFormat) DecodePattern(p *Pattern, reader io.Reader) error {
	// Read into the struct.
	err := binary.Read(reader, binary.LittleEndian, sf)
	if err != nil {
		return err
	}

	// Verify the magic header.
	if !sf.validMagic() {
		return fmt.Errorf("Magic header is wrong, got: %s", sf.Magic)
	}

	// Set fields from the header.
	p.Version = strings.Trim(string(sf.VersionBytes[:]), "\x00")
	p.Tempo = sf.Tempo

	return nil
}

// TrackBytes returns the number of bytes remaining for track data.
func (sf *sliceFormat) TrackBytes() int64 {
	return int64(sf.FileSize) - int64(len(sf.VersionBytes)) - 4 /* Tempo float32 */
}

// validMagic verifies that the data has the right kind of header.
func (sf sliceFormat) validMagic() bool {
	return string(sf.Magic[0:6]) == "SPLICE"
}

// trackFormat is the low level binary format for each track.
type trackFormat struct {
	ID       uint32
	NameSize byte
}

// DecodeTrack reads binary from reader and applies it to the Track.
func (tf *trackFormat) DecodeTrack(t *Track, reader io.Reader) error {

	// Decode header.
	err := binary.Read(reader, binary.LittleEndian, tf)
	if err != nil {
		return err
	}
	t.ID = int(tf.ID)

	// Decode name.
	name := make([]byte, tf.NameSize)
	_, err = io.ReadFull(reader, name)
	if err != nil {
		return err
	}
	t.Name = string(name)

	// Decode steps.
	var steps [16]byte
	err = binary.Read(reader, binary.LittleEndian, &steps)
	if err != nil {
		return err
	}
	for i, step := range steps {
		t.Steps[i] = step == 1
	}

	return nil
}
