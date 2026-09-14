package format

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/protobom/protobom/pkg/formats"
)

var (
	DefaultEncoding         = formats.JSON
	DefaultSPDXJSONVersion  = formats.SPDX23JSON
	DefaultSPDXTVVersion    = formats.SPDX23TV
	DefaultSPDX3Version     = formats.SPDX3JSON
	DefaultCycloneDXVersion = formats.CDX15JSON
	JSONFormatMap           = map[string]formats.Format{
		"spdx":     formats.SPDXFORMAT,
		"spdx-2.2": formats.SPDX22JSON,
		"spdx-2.3": formats.SPDX23JSON,

		"spdx3":      formats.SPDX3JSON,
		"spdx-3":     formats.SPDX3JSON,
		"spdx-3.0":   formats.SPDX3JSON,
		"spdx-3.0.1": formats.SPDX3JSON,

		"cyclonedx":     formats.CDXFORMAT,
		"cyclonedx-1.0": formats.CDX10JSON,
		"cyclonedx-1.1": formats.CDX11JSON,
		"cyclonedx-1.2": formats.CDX12JSON,
		"cyclonedx-1.3": formats.CDX13JSON,
		"cyclonedx-1.4": formats.CDX14JSON,
		"cyclonedx-1.5": formats.CDX15JSON,
		"cyclonedx-1.6": formats.CDX16JSON,
		"cyclonedx-1.7": formats.CDX17JSON,
	}

	TVFormatMap = map[string]formats.Format{
		"spdx":     formats.SPDXFORMAT,
		"spdx-2.2": formats.SPDX22TV,
		"spdx-2.3": formats.SPDX23TV,
	}

	XMLFormatMap = map[string]formats.Format{}

	JSONEncoding = formats.JSON
	TEXTEncoding = formats.TEXT
	SPDX         = formats.SPDXFORMAT
	CDX          = formats.CDXFORMAT

	EncodingMap = map[string]string{
		"json": formats.JSON,
		"xml":  formats.XML,
		"text": formats.TEXT,
	}
)

type Format struct {
	formats.Format
}

// Parse parses the format string into a formats.Format
func Parse(fs, encoding string) (*Format, error) {
	if fs == "" {
		return nil, errors.New("no format specified")
	}

	s := strings.ToLower(fs)
	var fm map[string]formats.Format

	switch encoding {
	case formats.JSON:
		fm = JSONFormatMap
	case formats.TEXT:
		fm = TVFormatMap
	case formats.XML:
		fm = XMLFormatMap
	default:
		return nil,
			fmt.Errorf("unknown encoding: %s", encoding)
	}
	f, ok := fm[s]
	if !ok {
		return nil, fmt.Errorf("unknown format: %s", fs)
	}

	if f == formats.SPDXFORMAT {
		if encoding == formats.JSON {
			return &Format{DefaultSPDXJSONVersion}, nil
		}

		if encoding == formats.TEXT {
			return &Format{DefaultSPDXTVVersion}, nil
		}
	}

	if f == formats.CDXFORMAT {
		return &Format{DefaultCycloneDXVersion}, nil
	}

	return &Format{f}, nil
}

func Detect(f io.ReadSeeker) (*Format, error) {
	s := formats.Sniffer{}
	format, err := s.SniffReader(f)
	if err != nil {
		return nil, fmt.Errorf("detecting SBOM format: %w", err)
	}

	if format == "" {
		return nil, nil
	}

	return &Format{format}, nil
}

// Inverse returns the format to convert to when none is specified. SPDX
// documents (2.x and 3.x) invert to CycloneDX; CycloneDX documents invert
// to SPDX 3 when JSON-encoded and to SPDX 2 tag-value otherwise.
func (f *Format) Inverse() (*Format, error) {
	switch f.Type() {
	case formats.SPDXFORMAT:
		return &Format{DefaultCycloneDXVersion}, nil
	case formats.CDXFORMAT:
		encoding := f.Encoding()
		if encoding == formats.JSON {
			return &Format{DefaultSPDX3Version}, nil
		}
		if encoding == formats.TEXT {
			return &Format{DefaultSPDXTVVersion}, nil
		}
	}

	return nil, errors.New("SBOM format unknown")
}

func (f *Format) String() string {
	return string(f.Format)
}
