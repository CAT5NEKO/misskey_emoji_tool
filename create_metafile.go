package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type JSONFile struct {
	MetaVersion int     `json:"metaVersion"`
	Host        string  `json:"host"`
	ExportedAt  string  `json:"exportedAt"`
	Emojis      []Emoji `json:"emojis"`
}

type Emoji struct {
	FileName   string      `json:"fileName"`
	Downloaded bool        `json:"downloaded"`
	Emoji      EmojiDetail `json:"emoji"`
}

type EmojiDetail struct {
	ID                                      string   `json:"id"`
	UpdatedAt                               string   `json:"updatedAt"`
	Name                                    string   `json:"name"`
	Host                                    string   `json:"host"`
	Category                                string   `json:"category"`
	OriginalURL                             string   `json:"originalUrl"`
	PublicURL                               string   `json:"publicUrl"`
	URI                                     string   `json:"uri"`
	Type                                    string   `json:"type"`
	Aliases                                 []string `json:"aliases"`
	License                                 string   `json:"license"`
	LocalOnly                               bool     `json:"localOnly"`
	IsSensitive                             bool     `json:"isSensitive"`
	RoleIDsThatCanBeUsedThisEmojiAsReaction []string `json:"roleIdsThatCanBeUsedThisEmojiAsReaction"`
}

type ConfigYAMLSchema struct {
	Host           string                `yaml:"host"`
	EmojiParameter ConfigYAMLEmojiSchema `yaml:"emojiParameter"`
}

type ConfigYAMLEmojiSchema struct {
	License     string `yaml:"license"`
	IsSensitive bool   `yaml:"isSensitive"`
	LocalOnly   bool   `yaml:"localOnly"`
	Category    string `yaml:"category"`
}

const (
	configYAML = "./cfg/config.yaml"
)

func MakeJSONFile(directory string) error {
	timeNow := time.Now().Format("2006-01-02T04:05:06Z")

	yamlFile, err := os.Open(configYAML)
	if err != nil {
		return fmt.Errorf("failed to open YAML file: %v", err)
	}
	defer yamlFile.Close()

	yamlContent, err := io.ReadAll(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %v", err)
	}

	var configYAMLStruct ConfigYAMLSchema
	if err := yaml.Unmarshal(yamlContent, &configYAMLStruct); err != nil {
		return fmt.Errorf("failed to unmarshal YAML content: %v", err)
	}

	jsonFile := &JSONFile{
		MetaVersion: 2,
		Host:        configYAMLStruct.Host,
		ExportedAt:  timeNow,
	}

	emojis := []Emoji{}
	err = filepath.WalkDir(directory, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ext := filepath.Ext(path); ext == ".png" || ext == ".PNG" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".GIF" {
			pattern := regexp.MustCompile(`.*\(\d+\).*`)
			if pattern.MatchString(path) {
				fmt.Println("duplicate file:", path)
				return nil
			}

			fileName := strings.ReplaceAll(filepath.Base(path), "-", "_")
			os.Rename(path, filepath.Dir(path)+"/"+fileName)

			emojiDetail := EmojiDetail{
				Name:        fileName[:len(filepath.Base(path))-len(filepath.Ext(path))],
				Category:    configYAMLStruct.EmojiParameter.Category,
				LocalOnly:   configYAMLStruct.EmojiParameter.LocalOnly,
				IsSensitive: configYAMLStruct.EmojiParameter.IsSensitive,
				License:     configYAMLStruct.EmojiParameter.License,
				Type:        "image/webp",
				Aliases:     []string{},
			}

			emoji := Emoji{
				FileName:   fileName,
				Downloaded: true,
				Emoji:      emojiDetail,
			}

			emojis = append(emojis, emoji)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory: %v", err)
	}

	jsonFile.Emojis = emojis
	jsonData, err := json.Marshal(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	var buffer bytes.Buffer
	if err := json.Indent(&buffer, jsonData, "", "  "); err != nil {
		return fmt.Errorf("failed to indent JSON: %v", err)
	}

	jsonFileName := directory + "/meta.json"
	jsonFileOutput, err := os.Create(jsonFileName)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %v", err)
	}
	defer jsonFileOutput.Close()

	if _, err := jsonFileOutput.Write(buffer.Bytes()); err != nil {
		return fmt.Errorf("failed to write JSON data: %v", err)
	}

	return nil
}
