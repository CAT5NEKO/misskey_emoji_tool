package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

func CheckNames(directory string) error {
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileName := filepath.Base(path)
			fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			snakeCaseRegex := regexp.MustCompile(`^[a-z0-9_]+$`)
			if !snakeCaseRegex.MatchString(fileNameWithoutExt) {
				fmt.Printf("💥 ファイル名 '%s' はアルファベットのスネークケースでないです。\n", fileName)
			}
			byteCount := utf8.RuneCountInString(fileNameWithoutExt)
			if byteCount > 15 {
				fmt.Printf("💥 ファイル名 '%s' は15字以上の名称になっており冗長です。\n", fileName)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
