package domain

import (
	"context"
	"io"

	"github.com/protobom/protobom/pkg/sbom"
)

type ConvertService interface {
	Convert(ctx context.Context, in io.ReadSeekCloser, out io.Writer) error
}

//go:generate mockgen -destination=mocks/mock_reader.go -package=mocks github.com/protobom/sbom-convert/pkg/convert Reader
type Reader interface {
	ParseStream(in io.ReadSeeker) (*sbom.Document, error)
}

//go:generate mockgen -destination=mocks/mock_writer.go -package=mocks github.com/protobom/sbom-convert/pkg/convert Writer
type Writer interface {
	WriteStream(doc *sbom.Document, out io.Writer) error
}
