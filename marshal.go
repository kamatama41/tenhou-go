package tenhou

import (
	"compress/gzip"
	"fmt"
	"io"
)

func Marshal(w io.Writer, m *MJLog, opts ...MarshalOption) error {
	options := &marshalOptions{
		withGzip: true,
	}
	for _, o := range opts {
		o(options)
	}

	if options.withGzip {
		gw := gzip.NewWriter(w)
		defer gw.Close()
		w = gw
	}
	e := &encoder{w}
	return e.marshal(m)
}

type encoder struct {
	w io.Writer
}

func (e *encoder) marshal(m *MJLog) error {
	if err := e.write(fmt.Sprintf(`<mjloggm ver="%s">`, m.Version)); err != nil {
		return err
	}
	// SHUFFLE
	if err := e.write(newXmlElement("SHUFFLE", "seed", m.Seed, "ref", "").Text()); err != nil {
		return err
	}
	// GO
	goElem, err := m.GameInfo.MarshalMJLog()
	if err != nil {
		return err
	}
	if err := e.write(goElem.Text()); err != nil {
		return err
	}
	// UN
	unElem, err := m.Players.MarshalMJLog()
	if err != nil {
		return err
	}
	if err := e.write(unElem.Text()); err != nil {
		return err
	}
	// TAIKYOKU
	if err := e.write(newXmlElement("TAIKYOKU", "oya", "0").Text()); err != nil {
		return err
	}
	// Events
	for _, event := range m.Events {
		elem, err := event.MarshalMJLog()
		if err != nil {
			return err
		}
		if err := e.write(elem.Text()); err != nil {
			return err
		}
	}
	return e.write("</mjloggm>")
}

func (e *encoder) write(txt string) error {
	bs := []byte(txt)
	written := 0
	for written < len(bs) {
		n, err := e.w.Write(bs[written:])
		if err != nil {
			return err
		}
		written += n
	}
	return nil
}

type marshalOptions struct {
	withGzip bool
}

type MarshalOption func(*marshalOptions)

func WithoutGzip() MarshalOption {
	return func(opts *marshalOptions) {
		opts.withGzip = false
	}
}
