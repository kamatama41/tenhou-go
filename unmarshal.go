package tenhou

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"regexp"
)

var (
	ErrInvalidElement = errors.New("invalid element")
	errUnknownElement = errors.New("unknown element")

	reTsumo *regexp.Regexp
	reDaHai *regexp.Regexp
)

func init() {
	reTsumo = regexp.MustCompile(`^([TUVW])[0-9]+$`)
	reDaHai = regexp.MustCompile(`^([DEFG])[0-9]+$`)
}

func Unmarshal(r io.Reader, opts ...UnmarshalOption) (*MJLog, error) {
	options := &unmarshalOptions{
		withRawXML: false,
	}
	for _, o := range opts {
		o(options)
	}

	if !options.withRawXML {
		gr, err := gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
		defer gr.Close()
		r = gr
	}
	xmlReader, err := newXmlReader(r)
	if err != nil {
		return nil, err
	}
	d := &decoder{xmlReader}
	return d.unmarshal()
}

type decoder struct {
	r *xmlReader
}

func (d *decoder) unmarshal() (*MJLog, error) {
	mjlog := &MJLog{}

	startTaikyoku := false
	for {
		e, err := d.r.next()
		if err != nil {
			if err == io.EOF {
				return mjlog, nil
			}
			return nil, err
		}

		if startTaikyoku {
			if err := d.unmarshalEvent(mjlog, e); err != nil {
				return nil, wrapInvalidElementErr(e, err)
			}
		} else {
			switch e.Name {
			case "mjloggm":
				mjlog.Version = e.AttrByName("ver")
			case "SHUFFLE":
				mjlog.Seed = e.AttrByName("seed")
			case "GO":
				if err := mjlog.GameInfo.UnmarshalMJLog(e); err != nil {
					return nil, wrapInvalidElementErr(e, err)
				}
			case "UN":
				// 対局開始前のUN
				if err := mjlog.Players.UnmarshalMJLog(e); err != nil {
					return nil, wrapInvalidElementErr(e, err)
				}
			case "TAIKYOKU":
				// oyaは常に0なので、保持しない
				startTaikyoku = true
			default:
				return nil, wrapInvalidElementErr(e, errUnknownElement)
			}
		}
	}
}

func (d *decoder) unmarshalEvent(m *MJLog, e XmlElement) error {
	var event MJLogElement
	if e.Name == "INIT" {
		event = &KyokuStart{}
	} else if reTsumo.MatchString(e.Name) {
		event = &Tsumo{}
	} else if reDaHai.MatchString(e.Name) {
		event = &Dahai{}
	} else if e.Name == "N" {
		event = &Naki{}
	} else if e.Name == "DORA" {
		event = &DoraOpen{}
	} else if e.Name == "REACH" {
		event = &Reach{}
	} else if e.Name == "BYE" {
		event = &Bye{}
	} else if e.Name == "UN" {
		event = &Reconnect{}
	} else if e.Name == "AGARI" {
		event = &Agari{}
	} else if e.Name == "RYUUKYOKU" {
		event = &Ryuukyoku{}
	} else {
		return errUnknownElement
	}
	if err := event.UnmarshalMJLog(e); err != nil {
		return err
	}
	m.Events = append(m.Events, event)
	return nil
}

func wrapInvalidElementErr(elem XmlElement, err error) error {
	return fmt.Errorf("%w on %s: %v", ErrInvalidElement, elem.Name, err)
}

func newUnknownAttrErr(name, value string) error {
	return fmt.Errorf("unknown attr %s %s", name, value)
}

type unmarshalOptions struct {
	withRawXML bool
}

type UnmarshalOption func(*unmarshalOptions)

func WithRawXML() UnmarshalOption {
	return func(opts *unmarshalOptions) {
		opts.withRawXML = true
	}
}
