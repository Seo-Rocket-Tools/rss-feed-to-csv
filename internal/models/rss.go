package models

import "encoding/xml"

// RSS represents the root RSS feed structure
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel represents the RSS channel containing items
type Channel struct {
	Items []Item `xml:"item"`
}

// Item represents a single RSS feed item
type Item struct {
	Title          string         `xml:"title"`
	Link           string         `xml:"link"`
	Description    string         `xml:"description"`
	PubDate        string         `xml:"pubDate"`
	ContentEncoded string         `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	MediaContent   []MediaContent `xml:"http://search.yahoo.com/mrss/ content"`
}

// MediaContent represents media content in RSS items
type MediaContent struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Medium string `xml:"medium,attr"`
}

// GetImageURL returns the first image URL found in the media content
func (item *Item) GetImageURL() string {
	for _, media := range item.MediaContent {
		if media.IsImage() {
			return media.URL
		}
	}
	return ""
}

// IsImage checks if the media content is an image
func (mc *MediaContent) IsImage() bool {
	if mc.Medium == "image" {
		return true
	}
	if mc.Type != "" && len(mc.Type) > 6 && mc.Type[:6] == "image/" {
		return true
	}
	return false
}