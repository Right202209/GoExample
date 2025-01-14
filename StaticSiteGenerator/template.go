package main

import (
	"fmt"
	"html/template"
	"strings"
)

func applyTemplate(templatePath string, data interface{}, cfg *Config) (string, error) {
	fmt.Printf("applyTemplate: templatePath: %s\n", templatePath)

	tmpl, err := template.ParseFiles(templatePath)

	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}
	fmt.Printf("applyTemplate: template parsed successfully\n")

	var buf strings.Builder

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}
	fmt.Printf("applyTemplate: template executed successfully\n")

	return buf.String(), nil
}
