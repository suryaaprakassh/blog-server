package posts 

import (
	"bytes"
	"io"
	"os"
	"strings"
	"github.com/yuin/goldmark"
)

type Post struct {
	Slug    string
	Title   string
	Content string
	Tags    []string
	Date    string
}

func (p *Post) parseConfig(line string) (string, string) {
	parts := strings.Split(line, ":")
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	return key, value
}

func (p *Post) creatTags(line string) []string {
	ret := make([]string, 0)
	for _, val := range strings.Split(line, ",") {
		ret = append(ret, strings.TrimSpace(val))
	}
	return ret
}

func (p *Post) GetMetaTags(content string) {
	it := strings.Lines(content)
	for line := range it {
		if strings.HasPrefix(line, "--") {
			if strings.TrimSpace(line) == "---" {
				break
			}
			line = strings.TrimPrefix(line, "--")
			key, value := p.parseConfig(line)
			switch key {
			case "title":
				p.Title = value
			case "date":
				p.Date = value
			case "slug":
				p.Slug = value
			case "tags":
				p.Tags = p.creatTags(value)
			default:
				panic("Unknown Keyword!")
			}
		} else {
			break
		}
	}

	var contentBuf bytes.Buffer
	for line := range it {
		contentBuf.Write([]byte(line))
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(contentBuf.Bytes(), &buf); err != nil {
		panic(err)
	}
	p.Content = string(buf.Bytes())
}


func NewPost(file *os.File) *Post {
	f := &Post{}
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	f.GetMetaTags(string(content))
	return f
}
