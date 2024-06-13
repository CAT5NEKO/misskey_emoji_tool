package main

import (
	"flag"
	"log"
)

func main() {
	makeJSON := flag.Bool("makejson", false, "一括で絵文字をインポートできるようにJSONファイルを作成します")
	checkName := flag.Bool("checkName", false, "絵文字の命名規則を検知します")
	flag.Parse()

	if *makeJSON {
		if flag.NArg() != 1 {
			log.Fatal("使用方法: ./emojiTool -makejson <ディレクトリパス>")
		}
		directory := flag.Arg(0)
		err := MakeJSONFile(directory)
		if err != nil {
			log.Fatal(err)
		}
	} else if *checkName {
		if flag.NArg() != 1 {
			log.Fatal("使用方法: ./emojiTool -checkName <ディレクトリパス>")
		}
		directory := flag.Arg(0)
		err := CheckNames(directory)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("使用方法: ./emojiTool [-makejson <ディレクトリパス> | -checkName" +
			" <ディレクトリパス>]")
	}
}
