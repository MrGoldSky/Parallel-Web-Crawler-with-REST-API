package parser

type PageData struct {
    Title string
    InternalLinks []string
    ExternalLinks []string
}

type Parser interface {
    Parse(html []byte) (PageData, error)
}