package tenhou

import (
	"strconv"
	"testing"
)

func TestGameInfo_MarshalMJLog_Roundtrip(t *testing.T) {
	for _, c := range []struct {
		typ  int
		name string
	}{
		{9, "四般南喰赤"},
		{137, "四上南喰赤"},
		{41, "四特南喰赤"},
		{169, "四鳳南喰赤"},
		{215, "三上東速"},
	} {
		t.Run(c.name, func(t *testing.T) {
			g := GameInfo{}
			e := newXmlElement("GO", "type", strconv.Itoa(c.typ), "lobby", "0")
			err := g.UnmarshalMJLog(e)
			if err != nil {
				t.Errorf("Failed to unmarshal %v", err)
			}

			if c.name != g.Name() {
				t.Errorf("Unexpected name expected:%+v actual:%+v", c.name, g.Name())
			}

			e, err = g.MarshalMJLog()
			if err != nil {
				t.Errorf("Failed to marshal %v", err)
			}
			if typ, _ := strconv.Atoi(e.AttrByName("type")); c.typ != typ {
				t.Errorf("Unexpected type expected:%+v actual:%+v", c.typ, typ)
			}
		})
	}
}
