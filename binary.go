package drum

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// sliceFormat is the low level binary format for the slice. After initializing
// this struct, you can use it to read or write Pattern objects into their
// binary format.
type sliceFormat struct {
	Magic        [13]byte
	FileSize     byte
	VersionBytes [32]byte
	Tempo        float32
}

// DecodePattern reads binary data from reader and returns a Pattern.
func (sf *sliceFormat) DecodePattern(reader io.Reader) (*Pattern, error) {
	p := &Pattern{}

	// Read into the struct.
	err := binary.Read(reader, binary.LittleEndian, sf)
	if err != nil {
		return p, err
	}

	// Verify the magic header.
	if !sf.validMagic() {
		return p, fmt.Errorf("Magic header is wrong, got: %s", sf.Magic)
	}

	// Set fields from the header.
	p.Version = strings.Trim(string(sf.VersionBytes[:]), "\x00")
	p.Tempo = sf.Tempo

	return p, nil
}

// TrackBytes returns the number of bytes remaining for track data.
func (sf *sliceFormat) TrackBytes() int64 {
	return int64(sf.FileSize) - int64(len(sf.VersionBytes)) - 4 /* Tempo float32 */
}

// EncodePattern takes data from the given pattern and stores it in this
// object. Afterwards, you can use Write to output that data.
func (sf *sliceFormat) EncodePattern(p *Pattern) error {
	sf.Magic = [13]byte{'S', 'P', 'L', 'I', 'C', 'E'}

	for i, c := range p.Version {
		sf.VersionBytes[i] = byte(c)
	}

	sf.Tempo = p.Tempo

	return nil
}

// SetFileSize updates the FileSize, accomodating for the current internal
// data, plus the given trackBytes. Call this after calculating how much track
// data is available, and before calling Write.
func (sf *sliceFormat) SetFileSize(trackBytes int64) {
	sf.FileSize = byte(trackBytes + int64(len(sf.VersionBytes)) + 4) /* Tempo float32 */
}

// Write outputs the binary slice format to the writer.
func (sf *sliceFormat) Write(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, sf); err != nil {
		return err
	}
	return nil
}

// validMagic verifies that the data has the right kind of header.
func (sf sliceFormat) validMagic() bool {
	return string(sf.Magic[0:6]) == "SPLICE"
}

// trackFormat is the low level binary format for each track.  After
// initializing this struct, you can use it to read or write Track objects into
// their binary format.
type trackFormat struct {
	trackHeader
	Name  string
	steps [16]byte
}

// DecodeTrack reads binary from reader and returns a Track.
func (tf *trackFormat) DecodeTrack(reader io.Reader) (*Track, error) {
	t := &Track{}

	// Decode header.
	err := binary.Read(reader, binary.LittleEndian, &tf.trackHeader)
	if err != nil {
		return t, err
	}
	t.ID = int(tf.ID)

	// Decode name.
	name := make([]byte, tf.NameSize)
	_, err = io.ReadFull(reader, name)
	if err != nil {
		return t, err
	}
	t.Name = string(name)

	// Decode steps.
	var steps [16]byte
	err = binary.Read(reader, binary.LittleEndian, &steps)
	if err != nil {
		return t, err
	}
	for i, step := range steps {
		t.Steps[i] = step == 1
	}

	return t, nil
}

// EncodeTrack stores the given Track data in this object. Afterwards, you can
// use Write to output that data.
func (tf *trackFormat) EncodeTrack(t *Track) error {
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

// Write outputs the binary track format to the writer.
func (tf *trackFormat) Write(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, tf.trackHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(w, tf.Name); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, tf.steps); err != nil {
		return err
	}
	return nil
}

// trackHeader is the low level binary format for the header of each track
type trackHeader struct {
	ID       uint32
	NameSize byte
}

// ByteSize is the total size of the track data.
func (th *trackHeader) ByteSize() int64 {
	// ID + NameSize + len(Name) + 16 steps
	return int64(4 + 1 + th.NameSize + 16)

}
