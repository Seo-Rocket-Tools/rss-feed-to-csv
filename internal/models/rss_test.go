package models

import (
	"encoding/xml"
	"testing"
)

func TestItem_GetImageURL(t *testing.T) {
	tests := []struct {
		name string
		item Item
		want string
	}{
		{
			name: "with image media content",
			item: Item{
				MediaContent: []MediaContent{
					{
						URL:    "https://example.com/image.jpg",
						Type:   "image/jpeg",
						Medium: "image",
					},
				},
			},
			want: "https://example.com/image.jpg",
		},
		{
			name: "with multiple media content",
			item: Item{
				MediaContent: []MediaContent{
					{
						URL:  "https://example.com/video.mp4",
						Type: "video/mp4",
					},
					{
						URL:  "https://example.com/image.png",
						Type: "image/png",
					},
				},
			},
			want: "https://example.com/image.png",
		},
		{
			name: "no image",
			item: Item{},
			want: "",
		},
		{
			name: "media without image type",
			item: Item{
				MediaContent: []MediaContent{
					{
						URL:  "https://example.com/audio.mp3",
						Type: "audio/mp3",
					},
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.GetImageURL(); got != tt.want {
				t.Errorf("GetImageURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMediaContent_IsImage(t *testing.T) {
	tests := []struct {
		name  string
		media MediaContent
		want  bool
	}{
		{
			name:  "medium is image",
			media: MediaContent{Medium: "image"},
			want:  true,
		},
		{
			name:  "type starts with image/",
			media: MediaContent{Type: "image/jpeg"},
			want:  true,
		},
		{
			name:  "type is image/png",
			media: MediaContent{Type: "image/png"},
			want:  true,
		},
		{
			name:  "not an image",
			media: MediaContent{Type: "video/mp4"},
			want:  false,
		},
		{
			name:  "empty media",
			media: MediaContent{},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.media.IsImage(); got != tt.want {
				t.Errorf("IsImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRSSUnmarshal(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<item>
			<title>Test Item</title>
			<link>https://example.com/item1</link>
			<description>Item Description</description>
			<pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>
		</item>
	</channel>
</rss>`

	var rss RSS
	err := xml.Unmarshal([]byte(xmlData), &rss)
	if err != nil {
		t.Fatalf("Failed to unmarshal RSS: %v", err)
	}

	if len(rss.Channel.Items) != 1 {
		t.Fatalf("len(Channel.Items) = %d, want 1", len(rss.Channel.Items))
	}

	item := rss.Channel.Items[0]
	if item.Title != "Test Item" {
		t.Errorf("Item.Title = %q, want %q", item.Title, "Test Item")
	}
	if item.Link != "https://example.com/item1" {
		t.Errorf("Item.Link = %q, want %q", item.Link, "https://example.com/item1")
	}
	if item.Description != "Item Description" {
		t.Errorf("Item.Description = %q, want %q", item.Description, "Item Description")
	}
}

func TestRSSWithMediaContent(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:media="http://search.yahoo.com/mrss/">
	<channel>
		<item>
			<title>Test Item</title>
			<media:content url="https://example.com/media.jpg" type="image/jpeg" />
		</item>
	</channel>
</rss>`

	var rss RSS
	err := xml.Unmarshal([]byte(xmlData), &rss)
	if err != nil {
		t.Fatalf("Failed to unmarshal RSS: %v", err)
	}

	if len(rss.Channel.Items) != 1 {
		t.Fatalf("len(Channel.Items) = %d, want 1", len(rss.Channel.Items))
	}

	item := rss.Channel.Items[0]
	if len(item.MediaContent) == 0 {
		t.Fatal("Item.MediaContent is empty")
	}

	if item.MediaContent[0].URL != "https://example.com/media.jpg" {
		t.Errorf("MediaContent[0].URL = %q, want %q", item.MediaContent[0].URL, "https://example.com/media.jpg")
	}
}