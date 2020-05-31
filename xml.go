package tenhou

import (
	"encoding/xml"
	"fmt"
	"io"
)

type XmlElement struct {
	Name string
	Attr []XmlAttr
}

type XmlAttr struct {
	Name  string
	Value string
}

func (e XmlElement) Text() string {
	attrStr := ""
	for _, attr := range e.Attr {
		attrStr += fmt.Sprintf(` %s="%s"`, attr.Name, attr.Value)
	}
	return fmt.Sprintf("<%s%s/>", e.Name, attrStr)
}

func (e *XmlElement) AppendAttr(name, value string) {
	e.Attr = append(e.Attr, XmlAttr{
		Name:  name,
		Value: value,
	})
}

func (e *XmlElement) AttrByName(name string) string {
	for _, attr := range e.Attr {
		if attr.Name == name {
			return attr.Value
		}
	}
	return ""
}

func (e *XmlElement) ForEachAttr(f func(name, value string) error) error {
	for _, attr := range e.Attr {
		if err := f(attr.Name, attr.Value); err != nil {
			return err
		}
	}
	return nil
}

func newXmlElement(name string, nameValue ...string) XmlElement {
	elem := XmlElement{
		Name: name,
	}
	for i := 0; i < len(nameValue); i += 2 {
		name := nameValue[i]
		value := nameValue[i+1]
		elem.Attr = append(elem.Attr, XmlAttr{
			Name:  name,
			Value: value,
		})
	}
	return elem
}

type xmlReader struct {
	dec *xml.Decoder
}

func newXmlReader(r io.Reader) (*xmlReader, error) {
	return &xmlReader{dec: xml.NewDecoder(r)}, nil
}

func (r *xmlReader) next() (XmlElement, error) {
	for {
		t, err := r.dec.Token()
		if err != nil {
			return XmlElement{}, err
		}
		switch t := t.(type) {
		case xml.StartElement:
			var attr []XmlAttr
			for _, a := range t.Attr {
				attr = append(attr, XmlAttr{
					Name:  a.Name.Local,
					Value: a.Value,
				})
			}
			return XmlElement{
				Name: t.Name.Local,
				Attr: attr,
			}, nil
		}
	}
}
