package convert

import (
	"context"
	"io"

	"github.com/bom-squad/protobom/pkg/reader"
	"github.com/bom-squad/protobom/pkg/writer"

	"github.com/bom-squad/go-cli/internal/domain"
	"github.com/bom-squad/go-cli/pkg/format"
)

var _ domain.ConvertService = (*Service)(nil)

type Service struct {
	Format *format.Format
}

func NewService(opts ...Option) *Service {
	s := &Service{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type Option func(s *Service)

func (s *Service) Convert(_ context.Context, in io.ReadSeekCloser, out io.WriteCloser) error {
	r := reader.New()
	doc, err := r.ParseStream(in)
	if err != nil {
		return err
	}

	w := writer.New()
	w.Options.Format = s.Format.Format

	return w.WriteStream(doc, out)
}
