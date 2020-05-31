package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kamatama41/tenhou-go"
)

var (
	path        = flag.String("path", "testdata/sample.mjlog", "mjlogファイルがあるディレクトリ or ファイルのパス")
	withRawXml  = flag.Bool("with-raw-xml", false, "インプットのファイルが生のxmlの場合使う")
	withoutGzip = flag.Bool("without-gzip", false, "アウトプットのファイルを生のxmlにしたい場合使う")
)

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "", 0)
	if err := filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".mjlog") {
			in, err := os.OpenFile(path, os.O_RDONLY, 0644)
			if err != nil {
				return err
			}
			defer in.Close()

			logger.Printf("### %s", info.Name())

			var uOpts []tenhou.UnmarshalOption
			if *withRawXml {
				uOpts = append(uOpts, tenhou.WithRawXML())
			}
			mjlog, err := tenhou.Unmarshal(in, uOpts...)
			if err != nil {
				return err
			}

			out, err := os.OpenFile(fmt.Sprintf("out/copy-%s", info.Name()), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
			defer out.Close()

			var mOpts []tenhou.MarshalOption
			if *withoutGzip {
				mOpts = append(mOpts, tenhou.WithoutGzip())
			}
			return tenhou.Marshal(out, mjlog, mOpts...)
		}
		return nil
	}); err != nil {
		panic(err)
	}
}
