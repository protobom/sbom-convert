package convert

import (
	"github.com/bom-squad/sbom-convert/pkg/domain"
	"github.com/bom-squad/sbom-convert/pkg/format"
)

func WithFormat(f *format.Format) Option {
	return func(s *Service) {
		s.Format = f
	}
}

func WithReader(r domain.Reader) Option {
	return func(s *Service) {
		s.r = r
	}
}

func WithWriter(w domain.Writer) Option {
	return func(s *Service) {
		s.w = w
	}
}
