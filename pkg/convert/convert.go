package convert

import (
	"context"
	"io"

	"github.com/bom-squad/protobom/pkg/reader"
	"github.com/bom-squad/protobom/pkg/writer"

	"github.com/bom-squad/sbom-convert/pkg/domain"
	"github.com/bom-squad/sbom-convert/pkg/format"
)

var _ domain.ConvertService = (*Service)(nil)

type Service struct {
	Format *format.Format
	r      domain.Reader
	w      domain.Writer
}

func NewService(opts ...Option) *Service {
	s := &Service{}

	for _, opt := range opts {
		opt(s)
	}

	if s.r == nil {
		s.r = reader.New()
	}

	if s.w == nil {
		s.w = writer.New(
			writer.WithFormat(s.Format.Format),
		)
	}

	return s
}

type Option func(s *Service)

func (s *Service) Convert(_ context.Context, in io.ReadSeekCloser, out io.WriteCloser) error {
	doc, err := s.r.ParseStream(in)
	if err != nil {
		return err
	}

	return s.w.WriteStream(doc, out)
}
