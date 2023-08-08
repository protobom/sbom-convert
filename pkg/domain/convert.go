package domain

import (
	"context"
	"io"
)

type ConvertService interface {
	Convert(ctx context.Context, in io.ReadSeekCloser, out io.WriteCloser) error
}
