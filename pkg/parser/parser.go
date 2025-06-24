package parser

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
    Title string
    InternalLinks []string
    ExternalLinks []string
}

type Parser interface {
    Parse(html []byte) (PageData, error)
}

type HTMLParser struct {
    base *url.URL
}

func NewParser(baseURL string) (Parser, error) {
    u, err := url.Parse(baseURL)
    if err != nil {
        return nil, err
    }
    return &HTMLParser{base: u}, nil
}

func (p *HTMLParser) Parse(html []byte) (PageData, error) {
    doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
    if err != nil {
        return PageData{}, err
    }
    data := PageData{}
    // Title
    data.Title = strings.TrimSpace(doc.Find("title").Text())
    // Links
    doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
        href, _ := s.Attr("href")
        u, err := p.base.Parse(strings.TrimSpace(href))
        if err != nil {
            return
        }
        str := u.String()
        if u.Host == p.base.Host {
            data.InternalLinks = append(data.InternalLinks, str)
        } else {
            data.ExternalLinks = append(data.ExternalLinks, str)
        }
    })
    return data, nil
}