package tenhou

import (
	"testing"
)

func TestMarshalMentsu_Roundtrip(t *testing.T) {
	tc := map[string]struct {
		m       int
		name    string
		haiType HaiType
	}{
		"67pから8をチー": {
			m:       39103,
			name:    "チー",
			haiType: HaiTypePin,
		},
		"下家から南をポン": {
			m:       44105,
			name:    "ポン",
			haiType: HaiTypeJi,
		},
		"対面からポンした3sを加槓": {
			m:       30738,
			name:    "加槓",
			haiType: HaiTypeSou,
		},
		"3pを暗槓": {
			m:       11776,
			name:    "暗槓",
			haiType: HaiTypePin,
		},
		"2mを明槓": {
			m:       1027,
			name:    "明槓",
			haiType: HaiTypeMan,
		},
		"抜きドラ": {
			m:       31264,
			name:    "抜きドラ",
			haiType: HaiTypeJi,
		},
	}

	for n, c := range tc {
		t.Run(n, func(t *testing.T) {
			mentsu, err := UnmarshalMentsu(c.m)
			if err != nil {
				t.Fatalf("Failed to unmarshal %d %v", c.m, err)
			}
			if c.name != mentsu.Name() {
				t.Errorf("Unexpected name\nexpected:%s actual:%s", c.name, mentsu.Name())
			}
			hai := mentsu.What()
			if c.haiType != hai.Type {
				t.Errorf("Unexpected haiType\nexpected:%d actual:%d", c.haiType, hai.Type)
			}

			m := MarshalMentsu(mentsu)
			if c.m != m {
				t.Errorf("Unexpected m %+v\nexpected:%d actual:%d", mentsu, c.m, m)
			}
		})
	}
}
