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
	p := &Pattern{}

	// Open the file.
	fi, err := os.Open(path)
	if err != nil {
		return p, err
	}
	defer fi.Close()

	// Parse the file header.
	header := header{}
	trackBytes, err := decodeHeader(&header, fi)
	if err != nil {
		return p, err
	}

	// Set fields from the header.
	p.Version = header.Version()
	p.Tempo = header.Tempo

	// Read the rest of the file for tracks.
	// NOTE: without a bufio.Reader we only get one track before EOF
	trackReader := bufio.NewReader(io.LimitReader(fi, trackBytes))

	// Parse the track data.
	for {
		t := Track{}
		err := decodeTrack(&t, trackReader)
		if err != nil {
			fmt.Printf("ERR Reading Track: %s\n", err)
			if err == io.EOF {
				break
			}
			return p, err
		}
		p.AddTrack(&t)
	}

	return p, nil
}

// header is the low level binary header.
type header struct {
	Magic        [13]byte
	FileSize     byte
	VersionBytes [32]byte
	Tempo        float32
}

// ValidMagic verifies if the file has the right kind of header.
func (h header) ValidMagic() bool {
	return string(h.Magic[0:6]) == "SPLICE"
}

// Version returns the file version string.
func (h header) Version() string {
	return strings.Trim(string(h.VersionBytes[:]), "\x00")
}

// decodeHeader extracts binary from reader into the header. It returns the
// number of bytes remaining in the file for track data.
func decodeHeader(h *header, reader io.Reader) (int64, error) {
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

// decodeTrack extracts binary from the reader into the Track.
func decodeTrack(t *Track, reader io.Reader) error {
	buf := bufio.NewReader(reader)

	// Set the ID
	fmt.Printf("Reading ID\n")
	err := binary.Read(reader, binary.LittleEndian, &t.ID)
	if err != nil {
		return err
	}

	// Get the size of the Name.
	fmt.Printf("Reading Size\n")
	size, err := buf.ReadByte()
	if err != nil {
		return err
	}

	// Set the Name.
	fmt.Printf("Reading Name\n")
	name := make([]byte, size)
	_, err = io.ReadFull(buf, name)
	if err != nil {
		return err
	}
	t.Name = string(name)

	// Set the Steps.
	fmt.Printf("Reading Steps\n")
	for i := 0; i < 16; i++ {
		b, err := buf.ReadByte()
		if err != nil {
			return err
		}
		t.Steps[i] = b == 1
	}

	return nil
}
