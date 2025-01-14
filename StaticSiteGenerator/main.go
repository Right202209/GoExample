// package main
//
// import (
//
//	"flag"
//	"fmt"
//	"html/template"
//	"io"
//	"log"
//	"os"
//	"path/filepath"
//	"strings"
//
// )
//
//	func main() {
//		defer func() {
//			if r := recover(); r != nil {
//				log.Fatalf("程序发生致命错误: %v", r)
//			}
//		}()
//
//		var (
//			contentDir  = flag.String("content", "content", "Content directory")
//			outputDir   = flag.String("output", "public", "Output directory")
//			templateDir = flag.String("template", "templates", "Template directory")
//			configFile  = flag.String("config", "config.yaml", "Configuration file")
//		)
//		flag.Parse()
//
//		fmt.Print(*templateDir)
//
//		// 读取配置文件
//		cfg, err := LoadConfig(*configFile)
//		if err != nil {
//			log.Fatalf("加载配置文件错误: %v", err)
//		}
//
//		if cfg.BaseURL == "" {
//			log.Fatal("配置错误: BaseURL 不能为空")
//		}
//
//		fmt.Printf("配置已加载: %+v\n", cfg)
//
//		// 读取标签信息
//		tags, err := LoadTags("data/tags.yaml")
//		if err != nil {
//			log.Printf("Error loading tags: %v", err)
//		}
//
//		fmt.Printf("Tags: %+v\n", tags)
//
//		// 创建输出目录
//		if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
//			log.Fatalf("Error creating output directory: %v", err)
//		}
//
//		// 创建 static 目录， 复制 static 目录的所有文件到输出目录
//		staticDir := filepath.Join("static")
//		err = filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
//			if err != nil {
//				return err
//			}
//
//			outputStaticPath := filepath.Join(*outputDir, strings.TrimPrefix(path, staticDir))
//			if info.IsDir() {
//				if err := os.MkdirAll(outputStaticPath, os.ModePerm); err != nil {
//					return fmt.Errorf("failed to create output static dir: %w", err)
//				}
//				return nil
//			}
//
//			if err := copyFile(path, outputStaticPath); err != nil {
//				return fmt.Errorf("failed to copy static file: %w", err)
//			}
//			return nil
//
//		})
//		if err != nil {
//			log.Fatalf("Error coping static dir: %v", err)
//		}
//
//		// 遍历内容目录下的所有Markdown文件
//		err = filepath.Walk(*contentDir, func(path string, info os.FileInfo, err error) error {
//			if err != nil {
//				return err
//			}
//
//			if info.IsDir() {
//				if info.Name() != "posts" {
//					return nil // 跳过目录
//				}
//
//				files, err := os.ReadDir(path)
//				if err != nil {
//					return fmt.Errorf("read content dir error: %w", err)
//				}
//
//				var posts []*Post
//				for _, file := range files {
//					if file.IsDir() {
//						continue
//					}
//
//					if strings.ToLower(filepath.Ext(file.Name())) != ".md" {
//						continue
//					}
//
//					mdPath := filepath.Join(path, file.Name())
//
//					post, err := buildPost(mdPath, tags, *outputDir, cfg)
//					if err != nil {
//						return fmt.Errorf("build post error: %w", err)
//					}
//					posts = append(posts, post)
//				}
//
//				err = renderListTemplate(posts, *templateDir, cfg, *outputDir)
//				if err != nil {
//					return fmt.Errorf("render list template error: %w", err)
//				}
//
//				return nil
//			}
//
//			if strings.ToLower(filepath.Ext(path)) != ".md" {
//				return nil
//			}
//			if info.Name() == "_index.md" { // 处理首页
//				// 读取 Markdown 文件
//				markdownContent, err := os.ReadFile(path)
//				if err != nil {
//					log.Printf("Error reading markdown file %s: %v\n", path, err)
//					return nil // 跳过当前文件
//				}
//
//				// 解析 Markdown
//				htmlContent, err := parseMarkdown(string(markdownContent))
//				if err != nil {
//					log.Printf("Error parsing markdown file %s: %v\n", path, err)
//					return nil
//				}
//
//				// 查找模板文件
//				templatePath := filepath.Join(*templateDir, "index.html")
//
//				// 应用模板
//				finalHTML, err := applyTemplate(templatePath, htmlContent, cfg)
//				if err != nil {
//					log.Printf("Error applying template for %s: %v\n", path, err)
//					return nil
//				}
//
//				// 生成输出文件路径，移除内容目录前缀，替换扩展名
//				outputFilePath := filepath.Join(*outputDir, strings.TrimPrefix(path, *contentDir))
//				outputFilePath = strings.Replace(outputFilePath, ".md", ".html", 1)
//
//				if err := os.WriteFile(outputFilePath, []byte(finalHTML), os.ModePerm); err != nil {
//					log.Printf("Error writing output HTML file %s: %v\n", outputFilePath, err)
//				}
//
//				log.Printf("Generated %s\n", outputFilePath)
//
//				return nil
//
//			}
//
//			post, err := buildPost(path, tags, *outputDir, cfg)
//
//			if err != nil {
//				return fmt.Errorf("build post error: %w", err)
//			}
//
//			if err := renderPostTemplate(post, *templateDir, cfg, *outputDir); err != nil {
//				return fmt.Errorf("render post template error: %w", err)
//			}
//
//			log.Printf("Generated %s\n", path)
//
//			return nil
//		})
//
//		if err != nil {
//			log.Fatalf("Error walking through content dir: %v", err)
//		}
//
//		fmt.Println("Static site generation complete!")
//	}
//
//	func copyFile(src, dst string) error {
//		sourceFile, err := os.Open(src)
//		if err != nil {
//			return fmt.Errorf("打开源文件错误: %w", err)
//		}
//		defer sourceFile.Close()
//
//		destFile, err := os.Create(dst)
//		if err != nil {
//			return fmt.Errorf("创建目标文件错误: %w", err)
//		}
//		defer destFile.Close()
//
//		_, err = io.Copy(destFile, sourceFile)
//		if err != nil {
//			return fmt.Errorf("复制文件错误: %w", err)
//		}
//		return nil
//	}
//
//	func buildPost(path string, tags map[string]Tag, outputDir string, cfg *Config) (*Post, error) {
//		log.Printf("开始处理文章: %s", path)
//		// 读取 Markdown 文件
//		markdownContent, err := os.ReadFile(path)
//		if err != nil {
//			return nil, fmt.Errorf("Error reading markdown file %s: %v", path, err)
//		}
//
//		// 解析 Markdown
//		htmlContent, err := parseMarkdown(string(markdownContent))
//		if err != nil {
//			return nil, fmt.Errorf("Error parsing markdown file %s: %v", path, err)
//		}
//
//		// 解析文章元数据
//		metadata, _ := parseMetadata(string(markdownContent))
//
//		var postTags []Tag
//		for _, t := range metadata.Tags {
//			if tag, ok := tags[t]; ok {
//				postTags = append(postTags, tag)
//			} else {
//				log.Printf("Tag %s not found", t)
//			}
//		}
//
//		// 生成输出文件路径，移除内容目录前缀，替换扩展名
//		outputFilePath := filepath.Join(outputDir, strings.TrimPrefix(path, "content"))
//		outputFilePath = strings.Replace(outputFilePath, ".md", ".html", 1)
//		post := &Post{
//			Title:       metadata.Title,
//			Date:        metadata.Date,
//			Content:     template.HTML(htmlContent),
//			Tags:        postTags,
//			OutputPath:  outputFilePath,
//			ContentDir:  strings.TrimPrefix(path, "content"),
//			BaseURL:     cfg.BaseURL,
//			Description: metadata.Description,
//		}
//
//		post.Link = post.BaseURL + "/" + strings.Replace(post.ContentDir, ".md", ".html", 1)
//
//		if len(postTags) == 0 {
//			log.Printf("警告: 文章 %s 没有标签", path)
//		}
//
//		log.Printf("文章处理完成: %s -> %s", path, outputFilePath)
//		return post, nil
//	}
//
//	func renderPostTemplate(post *Post, templateDir string, cfg *Config, outputDir string) error {
//		templatePath := filepath.Join(templateDir, "post.html")
//		finalHTML, err := applyTemplate(templatePath, post, cfg)
//		if err != nil {
//			return fmt.Errorf("Error applying template for %s: %v", post.ContentDir, err)
//		}
//
//		// 写出HTML文件
//		if err := os.WriteFile(post.OutputPath, []byte(finalHTML), os.ModePerm); err != nil {
//			return fmt.Errorf("Error writing output HTML file %s: %v", post.OutputPath, err)
//		}
//
//		return nil
//	}
//
//	func renderListTemplate(posts []*Post, templateDir string, cfg *Config, outputDir string) error {
//		templatePath := filepath.Join(templateDir, "list.html")
//		finalHTML, err := applyTemplate(templatePath, posts, cfg)
//		if err != nil {
//			return fmt.Errorf("Error applying template for list.html: %v", err)
//		}
//
//		outputFilePath := filepath.Join(outputDir, "posts", "index.html")
//
//		if err := os.WriteFile(outputFilePath, []byte(finalHTML), os.ModePerm); err != nil {
//			return fmt.Errorf("Error writing output HTML file: %v", outputFilePath, err)
//		}
//
//		return nil
//	}
package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("程序发生致命错误: %v", r)
		}
	}()

	var (
		contentDir  = flag.String("content", "content", "Content directory")
		outputDir   = flag.String("output", "public", "Output directory")
		templateDir = flag.String("template", "templates", "Template directory")
		configFile  = flag.String("config", "config.yaml", "Configuration file")
	)
	flag.Parse()

	fmt.Print(*templateDir)

	// 读取配置文件
	cfg, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("加载配置文件错误: %v", err)
	}

	if cfg.BaseURL == "" {
		log.Fatal("配置错误: BaseURL 不能为空")
	}

	fmt.Printf("配置已加载: %+v\n", cfg)

	// 读取标签信息
	tags, err := LoadTags("data/tags.yaml")
	if err != nil {
		log.Printf("Error loading tags: %v", err)
	}

	fmt.Printf("Tags: %+v\n", tags)

	// 创建输出目录
	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	// 创建 static 目录， 复制 static 目录的所有文件到输出目录
	staticDir := filepath.Join("static")
	err = filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		outputStaticPath := filepath.Join(*outputDir, strings.TrimPrefix(path, staticDir))
		if info.IsDir() {
			if err := os.MkdirAll(outputStaticPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create output static dir: %w", err)
			}
			return nil
		}

		if err := copyFile(path, outputStaticPath); err != nil {
			return fmt.Errorf("failed to copy static file: %w", err)
		}
		return nil

	})
	if err != nil {
		log.Fatalf("Error coping static dir: %v", err)
	}

	// 遍历内容目录下的所有Markdown文件
	err = filepath.Walk(*contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() != "posts" {
				return nil // 跳过目录
			}

			files, err := os.ReadDir(path)
			if err != nil {
				return fmt.Errorf("read content dir error: %w", err)
			}

			var posts []*Post
			for _, file := range files {
				if file.IsDir() {
					continue
				}

				if strings.ToLower(filepath.Ext(file.Name())) != ".md" {
					continue
				}

				mdPath := filepath.Join(path, file.Name())

				post, err := buildPost(mdPath, tags, *outputDir, cfg)
				if err != nil {
					return fmt.Errorf("build post error: %w", err)
				}
				posts = append(posts, post)
			}

			err = renderListTemplate(posts, *templateDir, cfg, *outputDir)
			if err != nil {
				return fmt.Errorf("render list template error: %w", err)
			}

			return nil
		}

		if strings.ToLower(filepath.Ext(path)) != ".md" {
			return nil
		}
		if info.Name() == "_index.md" { // 处理首页
			// 读取 Markdown 文件
			markdownContent, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading markdown file %s: %v\n", path, err)
				return nil // 跳过当前文件
			}

			// 解析 Markdown
			htmlContent, err := parseMarkdown(string(markdownContent))
			if err != nil {
				log.Printf("Error parsing markdown file %s: %v\n", path, err)
				return nil
			}

			// 查找模板文件
			templatePath := filepath.Join(*templateDir, "index.html")

			// 应用模板
			finalHTML, err := applyTemplate(templatePath, htmlContent, cfg)
			if err != nil {
				log.Printf("Error applying template for %s: %v\n", path, err)
				return nil
			}

			// 生成输出文件路径，移除内容目录前缀，替换扩展名
			outputFilePath := filepath.Join(*outputDir, strings.TrimPrefix(path, *contentDir))
			outputFilePath = strings.Replace(outputFilePath, ".md", ".html", 1)

			if err := os.WriteFile(outputFilePath, []byte(finalHTML), os.ModePerm); err != nil {
				log.Printf("Error writing output HTML file %s: %v\n", outputFilePath, err)
			}

			fmt.Printf("Generated %s\n", outputFilePath)

			return nil

		}

		post, err := buildPost(path, tags, *outputDir, cfg)

		if err != nil {
			return fmt.Errorf("build post error: %w", err)
		}

		if err := renderPostTemplate(post, *templateDir, cfg, *outputDir); err != nil {
			return fmt.Errorf("render post template error: %w", err)
		}

		fmt.Printf("Generated %s\n", path)

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through content dir: %v", err)
	}

	fmt.Println("Static site generation complete!")
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件错误: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件错误: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("复制文件错误: %w", err)
	}
	return nil
}

//func buildPost(path string, tags map[string]Tag, outputDir string, cfg *Config) (*Post, error) {
//	fmt.Printf("开始处理文章: %s\n", path)
//	// 读取 Markdown 文件
//	markdownContent, err := os.ReadFile(path)
//	if err != nil {
//		return nil, fmt.Errorf("Error reading markdown file %s: %v", path, err)
//	}
//
//	// 解析 Markdown
//	htmlContent, err := parseMarkdown(string(markdownContent))
//	if err != nil {
//		return nil, fmt.Errorf("Error parsing markdown file %s: %v", path, err)
//	}
//
//	// 解析文章元数据
//	metadata, _ := parseMetadata(string(markdownContent))
//
//	var postTags []Tag
//	for _, t := range metadata.Tags {
//		if tag, ok := tags[t]; ok {
//			postTags = append(postTags, tag)
//		} else {
//			log.Printf("Tag %s not found", t)
//		}
//	}
//
//	// 生成输出文件路径，移除内容目录前缀，替换扩展名
//	outputFilePath := filepath.Join(outputDir, strings.TrimPrefix(path, "content"))
//	outputFilePath = strings.Replace(outputFilePath, ".md", ".html", 1)
//	post := &Post{
//		Title:       metadata.Title,
//		Date:        metadata.Date,
//		Content:     template.HTML(htmlContent),
//		Tags:        postTags,
//		OutputPath:  outputFilePath,
//		ContentDir:  strings.TrimPrefix(path, "content"),
//		BaseURL:     cfg.BaseURL,
//		Description: metadata.Description,
//	}
//
//	post.Link = post.BaseURL + "/" + strings.Replace(post.ContentDir, ".md", ".html", 1)
//
//	if len(postTags) == 0 {
//		log.Printf("警告: 文章 %s 没有标签", path)
//	}
//
//	fmt.Printf("文章处理完成: %s -> %s\n", path, outputFilePath)
//	return post, nil
//}
//func renderPostTemplate(post *Post, templateDir string, cfg *Config, outputDir string) error {
//	templatePath := filepath.Join(templateDir, "post.html")
//	finalHTML, err := applyTemplate(templatePath, post, cfg)
//	if err != nil {
//		return fmt.Errorf("Error applying template for %s: %v", post.ContentDir, err)
//	}
//
//	// 写出HTML文件
//	if err := os.WriteFile(post.OutputPath, []byte(finalHTML), os.ModePerm); err != nil {
//		return fmt.Errorf("Error writing output HTML file %s: %v", post.OutputPath, err)
//	}
//
//	return nil
//}
//
//func renderListTemplate(posts []*Post, templateDir string, cfg *Config, outputDir string) error {
//	templatePath := filepath.Join(templateDir, "list.html")
//	finalHTML, err := applyTemplate(templatePath, posts, cfg)
//	if err != nil {
//		return fmt.Errorf("Error applying template for list.html: %v", err)
//	}
//
//	outputFilePath := filepath.Join(outputDir, "posts", "index.html")
//
//	if err := os.WriteFile(outputFilePath, []byte(finalHTML), os.ModePerm); err != nil {
//		return fmt.Errorf("Error writing output HTML file: %v", outputFilePath, err)
//	}
//
//	return nil
//}

func buildPost(path string, tags map[string]Tag, outputDir string, cfg *Config) (*Post, error) {
	fmt.Printf("buildPost: 开始处理文章: %s\n", path)
	// 读取 Markdown 文件
	markdownContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading markdown file %s: %v", path, err)
	}

	fmt.Printf("buildPost: Markdown 文件内容:\n%s\n", string(markdownContent))
	// 解析 Markdown
	htmlContent, err := parseMarkdown(string(markdownContent))
	if err != nil {
		return nil, fmt.Errorf("Error parsing markdown file %s: %v", path, err)
	}
	fmt.Printf("buildPost: Markdown 解析后的 HTML 内容:\n%s\n", htmlContent)

	// 解析文章元数据
	metadata, _ := parseMetadata(string(markdownContent))
	fmt.Printf("buildPost: 文章元数据: %+v\n", metadata)

	var postTags []Tag
	for _, t := range metadata.Tags {
		if tag, ok := tags[t]; ok {
			postTags = append(postTags, tag)
		} else {
			log.Printf("Tag %s not found", t)
		}
	}

	// 生成输出文件路径，移除内容目录前缀，替换扩展名
	outputFilePath := filepath.Join(outputDir, strings.TrimPrefix(path, "content"))
	outputFilePath = strings.Replace(outputFilePath, ".md", ".html", 1)
	post := &Post{
		Title:       metadata.Title,
		Date:        metadata.Date,
		Content:     template.HTML(htmlContent),
		Tags:        postTags,
		OutputPath:  outputFilePath,
		ContentDir:  strings.TrimPrefix(path, "content"),
		BaseURL:     cfg.BaseURL,
		Description: metadata.Description,
	}

	post.Link = post.BaseURL + "/" + strings.Replace(post.ContentDir, ".md", ".html", 1)

	if len(postTags) == 0 {
		log.Printf("警告: 文章 %s 没有标签", path)
	}
	fmt.Printf("buildPost: 生成的 Post 结构体: %+v\n", post)
	fmt.Printf("buildPost: 文章处理完成: %s -> %s\n", path, outputFilePath)
	return post, nil
}
func renderPostTemplate(post *Post, templateDir string, cfg *Config, outputDir string) error {
	fmt.Printf("renderPostTemplate: 开始渲染文章模板， post = %+v\n", post)
	templatePath := filepath.Join(templateDir, "post.html")
	fmt.Printf("renderPostTemplate: templatePath = %s\n", templatePath)
	finalHTML, err := applyTemplate(templatePath, post, cfg)
	if err != nil {
		return fmt.Errorf("Error applying template for %s: %v", post.ContentDir, err)
	}
	// 写出HTML文件
	fmt.Printf("renderPostTemplate: 生成的 HTML:\n%s\n", finalHTML)

	if err := os.WriteFile(post.OutputPath, []byte(finalHTML), os.ModePerm); err != nil {
		return fmt.Errorf("Error writing output HTML file %s: %v", post.OutputPath, err)
	}

	return nil
}

func renderListTemplate(posts []*Post, templateDir string, cfg *Config, outputDir string) error {
	fmt.Printf("renderListTemplate: 开始渲染列表模板, posts = %+v\n", posts)
	templatePath := filepath.Join(templateDir, "list.html")
	finalHTML, err := applyTemplate(templatePath, posts, cfg)
	if err != nil {
		return fmt.Errorf("Error applying template for list.html: %v", err)
	}
	fmt.Printf("renderListTemplate: 生成的 HTML:\n%s\n", finalHTML)

	outputFilePath := filepath.Join(outputDir, "posts", "index.html")

	if err := os.WriteFile(outputFilePath, []byte(finalHTML), os.ModePerm); err != nil {
		return fmt.Errorf("Error writing output HTML file: %v", outputFilePath, err)
	}

	return nil
}

// ... 其他代码 ...
