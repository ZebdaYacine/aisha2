package httpapi

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan/domain"
)

func TestArtisanMediaDTOUsesFrontendFieldNamesAndOmitsObjectKey(t *testing.T) {
	item := artisanMediaDTO(domain.Media{
		ID:               "media-1",
		MediaKind:        "IMAGE",
		ObjectKey:        "artisans/profile/media-1",
		OriginalFilename: "profile.jpg",
		MediaType:        "image/jpeg",
		SizeBytes:        42,
		SortOrder:        1,
		Visibility:       "PRIVATE",
		URL:              "http://localhost:9000/signed",
	})
	if item.ID != "media-1" || item.MediaKind != "IMAGE" || item.URL == "" {
		t.Fatalf("unexpected media DTO: %#v", item)
	}
	encoded, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	value := string(encoded)
	if !strings.Contains(value, `"mediaKind":"IMAGE"`) || strings.Contains(value, "objectKey") {
		t.Fatalf("unexpected serialized media DTO: %s", value)
	}
}
