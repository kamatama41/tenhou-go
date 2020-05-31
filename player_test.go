package tenhou

import (
	"testing"
)

func TestPlayerName_Encode_Roundtrip(t *testing.T) {
	for _, c := range []struct {
		name    string
		encoded string
	}{
		{"testuser1", "%74%65%73%74%75%73%65%72%31"},
		{"テストユーザー１", "%E3%83%86%E3%82%B9%E3%83%88%E3%83%A6%E3%83%BC%E3%82%B6%E3%83%BC%EF%BC%91"},
	} {
		t.Run(c.name, func(t *testing.T) {
			n := PlayerName(c.name)
			encoded := n.Encode()
			if encoded != c.encoded {
				t.Errorf("Unexpected encoded name expected:%s actual:%s", c.encoded, encoded)
			}

			name, err := DecodePlayerName(encoded)
			if err != nil {
				t.Fatalf("Failed to decode %s %v", n, err)
			}
			if string(name) != c.name {
				t.Errorf("Unexpected name expected:%s actual:%s", c.name, name)
			}
		})
	}
}
