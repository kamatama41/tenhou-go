package tenhou

import (
	"fmt"
	"sort"
	"strings"
)

var (
	haiNameList = [][]string{
		{"一", "二", "三", "四", "五", "六", "七", "八", "九"},
		{"①", "②", "③", "④", "⑤", "⑥", "⑦", "⑧", "⑨"},
		{"１", "２", "３", "４", "５", "６", "７", "８", "９"},
		{"東", "南", "西", "北", "白", "發", "中"},
	}
)

type Hai struct {
	Type   HaiType
	Number int // 0-8, 字牌は0-6
	Index  int // 4枚あるうちのどれか (赤牌などで区別が必要, 0が赤)
}

type HaiType int

const (
	HaiTypeMan HaiType = iota
	HaiTypePin
	HaiTypeSou
	HaiTypeJi
)

/**
 * 0-35:    萬子
 * 36-71:   筒子
 * 72-107:  索子
 * 108-135: 字牌
 */
func (h *Hai) Unmarshal(id int) error {
	if id < 0 || id > 136 {
		return fmt.Errorf("invalid Hai ID %d", id)
	}
	h.Type = HaiType(id / 36)
	h.Number = id % 36 / 4
	h.Index = id % 4
	return nil
}

func (h Hai) String() string {
	return h.Name()
}

/**
 * 赤牌にはrがつく e.g. r５
 */
func (h Hai) Name() string {
	aka := ""
	if h.IsAka() {
		aka = "r"
	}
	return fmt.Sprintf("%s%s", aka, haiNameList[h.Type][h.Number])
}

func (h Hai) ID() int {
	return int(h.Type)*36 + h.Number*4 + h.Index
}

func (h Hai) IsAka() bool {
	return h.Type != HaiTypeJi && h.Number == 4 && h.Index == 0
}

type HaiList []Hai

func (t HaiList) String() string {
	var names []string
	for _, h := range t {
		names = append(names, h.Name())
	}
	return strings.Join(names, "")
}

func (t HaiList) IDs() []int {
	var ret []int
	for _, hai := range t {
		ret = append(ret, hai.ID())
	}
	return ret
}

func (t HaiList) Sort() HaiList {
	var res HaiList
	res = append(res, t...)
	sort.Slice(res, func(i, j int) bool { return res[i].ID() < res[j].ID() })
	return res
}
