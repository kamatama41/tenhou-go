package tenhou

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPlayers_Marshal_Roundtrip(t *testing.T) {
	for name, c := range map[string]struct {
		input struct {
			n0, n1, n2, n3, dan, rate, sx string
		}
		expected Players
	}{
		"四麻": {
			struct {
				n0, n1, n2, n3, dan, rate, sx string
			}{
				n0:   "%74%65%73%74%75%73%65%72%31",
				n1:   "%74%65%73%74%75%73%65%72%32",
				n2:   "%74%65%73%74%75%73%65%72%33",
				n3:   "%74%65%73%74%75%73%65%72%34",
				dan:  "1,2,3,4",
				rate: "1500.12,1600.34,1700.56,1800.78",
				sx:   "F,M,F,M",
			},
			Players{
				{"testuser1", 1500.12, 1, "F"},
				{"testuser2", 1600.34, 2, "M"},
				{"testuser3", 1700.56, 3, "F"},
				{"testuser4", 1800.78, 4, "M"},
			},
		},
		"三麻": {
			struct {
				n0, n1, n2, n3, dan, rate, sx string
			}{
				n0:   "%E3%83%86%E3%82%B9%E3%83%88%E3%83%A6%E3%83%BC%E3%82%B6%E3%83%BC%EF%BC%91",
				n1:   "%E3%83%86%E3%82%B9%E3%83%88%E3%83%A6%E3%83%BC%E3%82%B6%E3%83%BC%EF%BC%92",
				n2:   "%E3%83%86%E3%82%B9%E3%83%88%E3%83%A6%E3%83%BC%E3%82%B6%E3%83%BC%EF%BC%93",
				n3:   "",
				dan:  "1,2,3,0",
				rate: "1500.12,1600.34,1700.56,1500.00",
				sx:   "F,M,F,C",
			},
			Players{
				{"テストユーザー１", 1500.12, 1, "F"},
				{"テストユーザー２", 1600.34, 2, "M"},
				{"テストユーザー３", 1700.56, 3, "F"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			elem := newXmlElement("UN",
				"n0", c.input.n0, "n1", c.input.n1, "n2", c.input.n2, "n3", c.input.n3,
				"dan", c.input.dan, "rate", c.input.rate, "sx", c.input.sx,
			)
			players := Players{}
			if err := players.UnmarshalMJLog(elem); err != nil {
				t.Fatalf("Failed to unmarshal %v", err)
			}

			if diff := cmp.Diff(c.expected, players); diff != "" {
				t.Errorf("Unexpected players (-want +got):\n%s", diff)
			}

			roundtrip, err := players.MarshalMJLog()
			if err != nil {
				t.Fatalf("Failed to marshal %v", err)
			}

			if diff := cmp.Diff(elem, roundtrip); diff != "" {
				t.Errorf("Unexpected XML element (-want +got):\n%s", diff)
			}
		})
	}
}
