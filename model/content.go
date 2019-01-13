package model

import (
	"github.com/SlyMarbo/rss"
	"github.com/indes/rssflow/tgraph"
	"strings"
)

type Content struct {
	SourceID     uint
	HashID       string `gorm:"primary_key"`
	RawID        string
	RawLink      string
	Title        string
	TelegraphUrl string
	EditTime
}

func getContentByFeedItem(source *Source, item *rss.Item) (Content, error) {
	tgpUrl := ""

	html := item.Content
	if html == "" {
		html = item.Summary
	}

	html = strings.Replace(html, "<![CDATA[", "", -1)
	html = strings.Replace(html, "]]>", "", -1)
	tgpUrl = tgraph.PublishItem(item.Title, html)

	var c = Content{
		Title:        item.Title,
		SourceID:     source.ID,
		RawID:        item.ID,
		HashID:       genHashID(source.Link, item.ID),
		TelegraphUrl: tgpUrl,
		RawLink:      item.Link,
	}

	return c, nil
}

func GenContentAndCheckByFeedItem(s *Source, item *rss.Item) (*Content, bool, error) {
	db := getConnect()
	defer db.Close()
	var (
		content   Content
		isBroaded bool
	)

	hashID := genHashID(s.Link, item.ID)
	db.Where("hash_id=?", hashID).First(&content)
	if content.HashID == "" {
		isBroaded = false
		content, _ = getContentByFeedItem(s, item)
		db.Save(&content)
	} else {
		isBroaded = true
	}

	return &content, isBroaded, nil
}