package convert_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/protobom/protobom/pkg/formats"
	"github.com/protobom/protobom/pkg/reader"
	"github.com/protobom/protobom/pkg/sbom"
	"github.com/protobom/protobom/pkg/writer"
	"go.uber.org/mock/gomock"

	"github.com/protobom/sbom-convert/pkg/convert"
	"github.com/protobom/sbom-convert/pkg/convert/mocks"
	"github.com/protobom/sbom-convert/pkg/format"
)

type readSeekerCloser struct {
	io.ReadSeeker
}

func (rsc *readSeekerCloser) Close() error {
	return nil
}

type writerCloser struct {
	io.Writer
}

func (wc *writerCloser) Close() error {
	return nil
}

func readJSON(t *testing.T, dir, path string) (res string) {
	t.Helper()

	//nolint:gosec
	if f, err := os.Open(filepath.Join(dir, path)); err != nil {
		t.Fatal(err)
	} else {
		defer f.Close() //nolint:errcheck

		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, f); err != nil {
			t.Fatal(err)
		}

		res = buf.String()
	}

	return res
}

func TestService_Convert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		in  io.ReadSeekCloser
		out io.WriteCloser
	}

	type mockFunc func(in io.ReadSeekCloser, out io.WriteCloser) (*mocks.MockReader, *mocks.MockWriter)

	tests := []struct {
		name    string
		format  format.Format
		args    args
		mocks   mockFunc
		want    string
		wantErr bool
	}{
		{
			name:   "happy: convert sbom",
			format: format.Format{Format: formats.CDX15JSON},
			args: args{
				in:  &readSeekerCloser{strings.NewReader(readJSON(t, "testdata", "spdx-2.3.json"))},
				out: &writerCloser{&bytes.Buffer{}},
			},
			mocks: func(in io.ReadSeekCloser, out io.WriteCloser) (*mocks.MockReader, *mocks.MockWriter) {
				bom := &sbom.Document{}
				mreader := mocks.NewMockReader(ctrl)
				mreader.EXPECT().ParseStream(in).Return(bom, nil)

				mwriter := mocks.NewMockWriter(ctrl)
				mwriter.EXPECT().WriteStream(bom, out).Return(nil)

				return mreader, mwriter
			},

			wantErr: false,
		},
		{
			name:   "sad: reader error",
			format: format.Format{Format: formats.CDX15JSON},
			args: args{
				in:  &readSeekerCloser{},
				out: &writerCloser{},
			},
			mocks: func(in io.ReadSeekCloser, _ io.WriteCloser) (*mocks.MockReader, *mocks.MockWriter) {
				mreader := mocks.NewMockReader(ctrl)
				mreader.EXPECT().ParseStream(in).Return(nil, io.EOF)

				mwriter := mocks.NewMockWriter(ctrl)

				return mreader, mwriter
			},
			wantErr: true,
		},
		{
			name:   "sad: writer error",
			format: format.Format{Format: formats.CDX15JSON},
			args: args{
				in:  &readSeekerCloser{strings.NewReader(readJSON(t, "testdata", "spdx-2.3.json"))},
				out: &writerCloser{},
			},
			mocks: func(in io.ReadSeekCloser, out io.WriteCloser) (*mocks.MockReader, *mocks.MockWriter) {
				bom := &sbom.Document{}
				mreader := mocks.NewMockReader(ctrl)
				mreader.EXPECT().ParseStream(in).Return(bom, nil)

				mwriter := mocks.NewMockWriter(ctrl)
				mwriter.EXPECT().WriteStream(bom, out).Return(io.EOF)

				return mreader, mwriter
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt

			mreader, mwriter := tt.mocks(tt.args.in, tt.args.out)

			s := convert.NewService(
				convert.WithFormat(&tt.format),
				convert.WithReader(mreader),
				convert.WithWriter(mwriter),
			)

			err := s.Convert(context.Background(), tt.args.in, tt.args.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewService(t *testing.T) {
	type args struct {
		opts []convert.Option
	}

	treader := reader.New()
	twriter := writer.New()
	tests := []struct {
		name string
		args args
		want *convert.Service
	}{
		{
			name: "happy: new default service",
			args: args{
				opts: []convert.Option{
					convert.WithFormat(&format.Format{Format: formats.CDX15JSON}),
					convert.WithReader(nil),
					convert.WithWriter(nil),
				},
			},
			want: &convert.Service{
				Format: &format.Format{Format: formats.CDX15JSON},
			},
		},
		{
			name: "happy: new service with reader and writer",
			args: args{
				opts: []convert.Option{
					convert.WithFormat(&format.Format{Format: formats.CDX15JSON}),
					convert.WithReader(treader),
					convert.WithWriter(twriter),
				},
			},
			want: &convert.Service{
				Format: &format.Format{Format: formats.CDX15JSON},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convert.NewService(tt.args.opts...); !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(convert.Service{})) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}
