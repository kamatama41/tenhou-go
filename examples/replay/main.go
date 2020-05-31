package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kamatama41/tenhou-go"
)

var (
	path = flag.String("path", "testdata/sample.mjlog", "mjlogファイルがあるディレクトリ or ファイルのパス")
)

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "", 0)
	if err := filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".mjlog") {
			f, err := os.OpenFile(path, os.O_RDONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			logger.Printf("### %s", info.Name())
			mjlog, err := tenhou.Unmarshal(f)
			if err != nil {
				return err
			}

			tenhou.NewReader(mjlog, logger).ReadAll()
			logger.Printf("")
		}
		return nil
	}); err != nil {
		panic(err)
	}
}
