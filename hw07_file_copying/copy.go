package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	writer, err := os.Create(toPath)
	if err != nil {
		return err
	}

	reader, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	readerStat, err := reader.Stat()
	if err != nil {
		return err
	}

	if offset > readerStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	offset, err = reader.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	if limit <= 0 {
		limit = readerStat.Size()
	}

	if limit > readerStat.Size()-offset {
		limit = readerStat.Size() - offset
	}
	bar := pb.Full.Start64(limit)
	limitReader := io.LimitReader(reader, limit)
	barReader := bar.NewProxyReader(limitReader)

	for {
		readChunk, err := io.CopyN(writer, barReader, 1024)
		offset += readChunk
		if err == io.EOF {
			break
		}
	}

	bar.Finish()

	return nil
}
