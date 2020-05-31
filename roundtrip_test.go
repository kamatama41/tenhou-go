package tenhou

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRoundtrip(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("The test did panic. %v", r)
		}
	}()

	root := "testdata/sample.mjlog"

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		f, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		origin, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		// Try unmarshal
		mjlog, err := Unmarshal(bytes.NewReader(origin))
		if err != nil {
			return err
		}

		// Try Marshal
		copied := bytes.NewBuffer(nil)
		if err := Marshal(copied, mjlog); err != nil {
			return err
		}

		// Decode to raw XML
		xmlOrigin, err := decodeMJLogToXml(bytes.NewReader(origin))
		if err != nil {
			return err
		}
		xmlCopied, err := decodeMJLogToXml(bytes.NewReader(copied.Bytes()))
		if err != nil {
			return err
		}

		if diff := cmp.Diff(xmlOrigin, xmlCopied); diff != "" {
			t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
		}

		return nil
	}); err != nil {
		t.Fatalf("Failed %v", err)
	}
}

func decodeMJLogToXml(r io.Reader) (string, error) {
	r, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}
	xml, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(xml), nil
}
