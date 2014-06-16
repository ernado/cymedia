package cymedia

import (
	"io"
)

type Progressor struct {
	Length   int64
	Rate     int64
	Reader   *io.PipeReader
	Writer   *io.PipeWriter
	Progress chan float32
}

func (p *Progressor) Start() {
	defer p.Writer.Close()
	var total int64
	bufLen := p.Length * 1. / p.Rate
	p.Progress <- float32(0)
	for {
		buffer := make([]byte, bufLen)
		read, err := p.Reader.Read(buffer)
		if err == io.EOF {
			break
		}
		total += int64(read)
		if total == p.Length {
			break
		}
		p.Progress <- float32(total) / float32(p.Length) * 100
	}
	close(p.Progress)
}

func Progress(f io.Reader, length int64, rate int64, progress chan float32) io.Reader {
	progressReader, progressWriter := io.Pipe()
	reader := io.TeeReader(f, progressWriter)
	p := Progressor{length, rate, progressReader, progressWriter, progress}
	go p.Start()
	return reader
}
