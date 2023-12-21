package convert

import (
	"github.com/bom-squad/sbom-convert/pkg/format"
)

func WithFormat(f *format.Format) Option {
	return func(s *Service) {
		s.Format = f
	}
}
