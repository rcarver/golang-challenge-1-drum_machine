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

	// Parse the file header.
	var header sliceHeader
	trackBytes, err := decodeSlice(&header, buffer)
	if err != nil {
		return p, err
	}

	// Set fields from the header.
	p.Version = header.Version()
	p.Tempo = header.Tempo

	// Read the rest of the file for tracks.
	trackReader := io.LimitReader(buffer, trackBytes)

	// Parse the track data.
	for {
		t := Track{}
		err := decodeTrack(&t, trackReader)
		if err != nil {
			if err == io.EOF {
				break
			}
			return p, err
		}
		p.AddTrack(&t)
	}

	return p, nil
}

// sliceHeader is the low level binary header for the entire slice.
type sliceHeader struct {
	Magic        [13]byte
	FileSize     byte
	VersionBytes [32]byte
	Tempo        float32
}

// ValidMagic verifies if the file has the right kind of header.
func (h sliceHeader) ValidMagic() bool {
	return string(h.Magic[0:6]) == "SPLICE"
}

// Version returns the file version string.
func (h sliceHeader) Version() string {
	return strings.Trim(string(h.VersionBytes[:]), "\x00")
}

// decodeSlice extracts binary from reader into the header. It returns the
// number of bytes remaining in the file for track data.
func decodeSlice(h *sliceHeader, reader io.Reader) (int64, error) {
	// Read into the struct.
	err := binary.Read(reader, binary.LittleEndian, h)
	if err != nil {
		return 0, err
	}

	// Verify the magic header.
	if !h.ValidMagic() {
		return 0, fmt.Errorf("Magic header is wrong, got: %s", h.Magic)
	}

	// Calculate the remaining bytes after reading everything.
	bytes := int64(h.FileSize) - int64(len(h.VersionBytes)) - 4 /* Tempo float 32 */

	return bytes, nil
}

// trackHeader is the low level binary header for each track.
type trackHeader struct {
	ID       uint32
	NameSize byte
}

// decodeTrack extracts binary from the reader into the Track.
func decodeTrack(t *Track, reader io.Reader) error {
	var header trackHeader

	// Decode header.
	err := binary.Read(reader, binary.LittleEndian, &header)
	if err != nil {
		return err
	}
	t.ID = int(header.ID)

	// Decode name.
	name := make([]byte, header.NameSize)
	_, err = io.ReadFull(reader, name)
	if err != nil {
		return err
	}
	t.Name = string(name)

	// Decode steps.
	var steps [16]byte
	err = binary.Read(reader, binary.LittleEndian, &steps)
	for i, step := range steps {
		t.Steps[i] = step == 1
	}

	return nil
}
