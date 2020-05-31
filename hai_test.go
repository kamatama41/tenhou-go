package tenhou

import (
	"fmt"
	"reflect"
	"testing"
)

func TestHai_Unmarshal(t *testing.T) {
	for _, c := range []struct {
		id           int
		expected     Hai
		expectedName string
	}{
		{0, Hai{0, 0, 0}, "一"},
		{1, Hai{0, 0, 1}, "一"},
		{2, Hai{0, 0, 2}, "一"},
		{3, Hai{0, 0, 3}, "一"},
		{4, Hai{0, 1, 0}, "二"},
		{5, Hai{0, 1, 1}, "二"},
		{6, Hai{0, 1, 2}, "二"},
		{7, Hai{0, 1, 3}, "二"},
		{36, Hai{1, 0, 0}, "①"},
		{72, Hai{2, 0, 0}, "１"},
		{108, Hai{3, 0, 0}, "東"},
	} {
		t.Run(fmt.Sprintf("Unmarshal(%d)", c.id), func(t *testing.T) {
			hai := Hai{}
			err := hai.Unmarshal(c.id)
			if err != nil {
				t.Errorf("Failed to unmarshal %v", err)
			}
			if !reflect.DeepEqual(c.expected, hai) {
				t.Errorf("Unexpected Hai expected:%+v actual:%+v", c.expected, hai)
			}
			if c.expectedName != hai.Name() {
				t.Errorf("Unexpected name expected:%s actual:%s", c.expectedName, hai.Name())
			}
		})
	}
}
