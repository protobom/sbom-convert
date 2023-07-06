package translate

import (
	"errors"
	"fmt"
	"io"

	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/bom-squad/protobom/pkg/reader"
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/bom-squad/protobom/pkg/writer"
	"github.com/sirupsen/logrus"
)

type TranslatorImplementation interface {
	IngestSBOM(*TranslationOptions, io.ReadSeekCloser) (*sbom.Document, error)
	SerializeSBOM(*TranslationOptions, *sbom.Document, io.WriteCloser) error
	DetectFormat(io.ReadSeeker) (*formats.Format, error)
	ComputeTranslationFormat(*TranslationOptions, io.ReadSeeker) error
}

type defaultTranslatorImplementation struct{}

// IngestSBOM reads SBOM data from a reader and returns a protobom document.
// It will autodetect the format being ingested.
func (di *defaultTranslatorImplementation) IngestSBOM(_ *TranslationOptions, in io.ReadSeekCloser) (*sbom.Document, error) {
	r := reader.New()
	doc, err := r.ParseStream(in)
	if err != nil {
		return nil, fmt.Errorf("parsing sbom: %w", err)
	}
	return doc, nil
}

// SerializeSBOM writes a protobom to the stream in the format specified by
// the options
func (di *defaultTranslatorImplementation) SerializeSBOM(opts *TranslationOptions, doc *sbom.Document, out io.WriteCloser) error {
	if opts.Format == "" {
		return errors.New("no destination format defined")
	}

	// Create new writer
	w := writer.New()

	// Set the format
	w.Options.Format = opts.Format

	// Render the new serialized SBOM to the stream
	if err := w.WriteStream(doc, out); err != nil {
		return fmt.Errorf("writing sbom: %w", err)
	}
	return nil
}

// DetectFormat checks the in stream and returns the format of the SBOM
func (di *defaultTranslatorImplementation) DetectFormat(in io.ReadSeeker) (*formats.Format, error) {
	s := formats.Sniffer{}
	format, err := s.SniffReader(in)
	if err != nil {
		return nil, fmt.Errorf("detecting SBOM format: %w", err)
	}

	if format == "" {
		return nil, nil
	}

	return &format, nil
}

func (di *defaultTranslatorImplementation) ComputeTranslationFormat(o *TranslationOptions, in io.ReadSeeker) error {
	// If the options have a format defined, we use it
	if o.Format != "" {
		logrus.Debugf("Translating to specified format %s", o.Format)
		return nil
	}

	// If not, we need to know the original format...
	originalFormat, err := di.DetectFormat(in)
	if err != nil {
		return fmt.Errorf("detecting SBOM format: %w", err)
	}

	// To compute the inverse format ...
	newFormat, err := o.InverseFormat(originalFormat)
	if err != nil {
		return fmt.Errorf("unable to determine format automatically: %w", err)
	}

	logrus.Debugf("computing inverse format: %s", newFormat)
	o.Format = newFormat
	return nil
}
