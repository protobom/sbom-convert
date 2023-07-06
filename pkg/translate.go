package translate

import (
	"errors"
	"fmt"
	"io"

	"github.com/bom-squad/protobom/pkg/formats"
)

const (
	DefaultEncoding         = formats.JSON
	DefaultSPDXVersion      = formats.SPDX23JSON
	DefaultCycloneDXVersion = formats.CDX14JSON
)

func NewTranslator() *Translator {
	return &Translator{
		impl: &defaultTranslatorImplementation{},
	}
}

type Translator struct {
	impl TranslatorImplementation
}

type TranslationOptions struct {
	Format formats.Format
}

// InverseFormat gets am SBOM format and returns the default inverse format.
// For example, if the SBOM is a CycloneDX document, it will return the
// SPDX version marked as default.
func (to *TranslationOptions) InverseFormat(f *formats.Format) (formats.Format, error) {
	switch f.Type() {
	case formats.SPDXFORMAT:
		return DefaultCycloneDXVersion, nil
	case formats.CDXFORMAT:
		return DefaultSPDXVersion, nil
	default:
		return "", errors.New("SBOM format unknown")
	}
}

// defaultTranslationOptions are the options used when calling Translate()
// Ny default we will auto compute the target format based on the SBOM original
// format and the version we const'ed as default.
var defaultTranslationOptions = &TranslationOptions{
	Format: "",
}

// Translate reads an SBOM from the in stream, translates it using the default
func (t *Translator) Translate(in io.ReadSeekCloser, out io.WriteCloser) error {
	return t.TranslateWithOptions(defaultTranslationOptions, in, out)
}

// DetectFormat analyzes a stream and returns and if it is a known SBOM
// format returns a protobom format string.
func (t *Translator) DetectFormat(in io.ReadSeeker) (*formats.Format, error) {
	return t.impl.DetectFormat(in)
}

// TranslateWithOptions reads an SBOM from stream in and translates it applying
// options o.
func (t *Translator) TranslateWithOptions(o *TranslationOptions, in io.ReadSeekCloser, out io.WriteCloser) error {
	// Since we are modifying the options, copy them here
	opts := *o

	// Compute the translation format from the specified options
	if err := t.impl.ComputeTranslationFormat(&opts, in); err != nil {
		return fmt.Errorf("computing detination format")
	}

	// Ingest the SBOM from the stream
	doc, err := t.impl.IngestSBOM(&opts, in)
	if err != nil {
		return fmt.Errorf("ingesting SBOM: %w", err)
	}

	// Serialize it to the out stream in the new format
	if err := t.impl.SerializeSBOM(&opts, doc, out); err != nil {
		return fmt.Errorf("serializing SBOM: %w", err)
	}

	return nil
}
