package xmltransformer

import (
	"encoding/xml"
)

// RSS は main に定義されている構造体に準拠したXMLデータのルート要素です。
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel はRSS内のチャネル情報を表します。
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

// Item はRSS内の各記事情報を表します。
type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

// Transformer はXML変換のためのインターフェースです。
type Transformer interface {
	TransformToRSS(body []byte) (*RSS, error)
}

// XmlTransformer は Transformer インターフェースの実装です。
type XmlTransformer struct{}

// New は XmlTransformer の新しいインスタンスを生成します。
func New() *XmlTransformer {
	return &XmlTransformer{}
}

// TransformToRSS は入力されたXML文字列をパースし、RSS構造体に変換します。
func (t *XmlTransformer) TransformToRSS(body []byte) (*RSS, error) {
	var rss RSS
	err := xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}
	return &rss, nil
}
