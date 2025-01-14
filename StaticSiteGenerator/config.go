package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"os"
	"regexp"
)

// Config 配置文件
type Config struct {
	SiteTitle   string `yaml:"siteTitle"`
	Description string `yaml:"description"`
	BaseURL     string `yaml:"baseURL"`
	Theme       string `yaml:"theme"`
}
type Tag struct {
	Name        string `yaml:"name"`
	Slug        string `yaml:"slug"`
	Description string `yaml:"description"`
}
type Tags map[string]Tag

// Post 文章数据
type Post struct {
	Title       string
	Date        string
	Content     template.HTML
	Tags        []Tag
	OutputPath  string
	ContentDir  string
	Link        string
	BaseURL     string
	Description string
}
type Metadata struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Tags        []string `yaml:"tags"`
	Description string   `yaml:"description"`
}

// LoadConfig 读取配置文件
func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func LoadTags(path string) (map[string]Tag, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read tags error %w", err)
	}
	var tags map[string]Tag
	err = yaml.Unmarshal(file, &tags)
	if err != nil {
		return nil, fmt.Errorf("unmarshal tags error %w", err)
	}
	return tags, nil
}
func parseMetadata(content string) (Metadata, error) {
	metadata := Metadata{}
	re := regexp.MustCompile(`(?s)(---)(.+?)---`)
	match := re.FindStringSubmatch(content)

	if len(match) == 0 {
		return metadata, nil
	}
	yamlContent := match[2]
	err := yaml.Unmarshal([]byte(yamlContent), &metadata)
	if err != nil {
		return Metadata{}, fmt.Errorf("yaml unmarshal error %w", err)
	}
	return metadata, nil

}
