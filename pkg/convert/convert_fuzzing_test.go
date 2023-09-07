package convert

import (
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/bom-squad/protobom/pkg/reader"
)

type ReadWriteSeeker struct {
	*os.File
}

func (rws *ReadWriteSeeker) Close() error {
	return rws.File.Close()
}

func writeStringToTempFile(content string) (io.ReadSeekCloser, error) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "tempfile")
	if err != nil {
		return nil, err
	}

	// Write the content to the temporary file
	_, err = tempFile.WriteString(content)
	if err != nil {
		return nil, err
	}

	// Close the file to make sure all data is flushed to disk
	err = tempFile.Close()
	if err != nil {
		return nil, err
	}

	// Reopen the temporary file for reading and seeking
	file, err := os.OpenFile(tempFile.Name(), os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &ReadWriteSeeker{file}, nil
}

func ParseStreamWrapper(content string) {
	t := io.NopCloser(strings.NewReader(content))
	r := reader.New()
	t, _ = writeStringToTempFile(content)
	t2 := t.(io.ReadSeekCloser)
	r.ParseStream(t2)
}

func FuzzParseStream(f *testing.F) {
	// filePaths := []string{"/home/wei/code/sboms/python/abhiTronix/vidgear/syft_spdx.json", "/home/wei/code/sboms/python/awslabs/aws-data-wrangler/syft_cyclonedx.json"}
	filePaths := []string{"/home/wei/code/sboms/python/abhiTronix/vidgear/syft_spdx.json"}
	for _, filePath := range filePaths {
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(string(content))
		f.Add(string(content))
	}

	f.Fuzz(func(t *testing.T, orig string) {
		ParseStreamWrapper(orig)
	})
}
