package tenhou

import (
	"strconv"
)

type GameInfo struct {
	Demo     bool
	Akadora  bool
	Ariari   bool
	Tonpu    bool
	Sanma    bool
	Sokutaku bool
	Taku     Taku
	Lobby    string
}

func (g *GameInfo) UnmarshalMJLog(e XmlElement) error {
	return e.ForEachAttr(func(name, value string) error {
		switch name {
		case "lobby":
			// 大会のロビーID
			g.Lobby = value
		case "type":
			gameType, err := strconv.Atoi(value)
			if err != nil {
				return err
			}

			// 0x01 (0b00000001): 人間との対戦のとき 1
			g.Demo = (gameType & 0x01) == 0
			// 0x02 (0b00000010): 赤牌なしのとき 1
			g.Akadora = (gameType&0x02)>>1 == 0
			// 0x04 (0b00000100): ナシナシのとき 1
			g.Ariari = (gameType&0x04)>>2 == 0
			// 0x08 (0b00001000): 東風戦のとき0, 東南戦のとき 1
			g.Tonpu = (gameType&0x08)>>3 == 0
			// 0x10 (0b00010000): 三麻のとき 1
			g.Sanma = (gameType&0x10)>>4 == 1
			// 0x40 (0b01000000): 速卓(持ち時間が短い)のとき 1
			g.Sokutaku = (gameType&0x40)>>6 == 1

			// 8ビット中の最初の1ビットと3ビットの組み合わせで決まる
			// 0x00 (0b00000000): 一般
			// 0x80 (0b10000000): 上級
			// 0x20 (0b00100000): 特上
			// 0xA0 (0b10100000): 鳳凰
			g.Taku = Taku(gameType & 0xA0)
		default:
			return newUnknownAttrErr(name, value)
		}
		return nil
	})
}

func (g *GameInfo) MarshalMJLog() (XmlElement, error) {
	var result int
	// demoじゃないとき1
	if !g.Demo {
		result = result | 1
	}
	// 赤牌が無いとき1
	if !g.Akadora {
		result = result | 1<<1
	}
	// ナシナシのとき1
	if !g.Ariari {
		result = result | 1<<2
	}
	// 東南戦のとき1
	if !g.Tonpu {
		result = result | 1<<3
	}
	// サンマのとき1
	if g.Sanma {
		result = result | 1<<4
	}
	// 速卓のとき1
	if g.Sokutaku {
		result = result | 1<<6
	}
	result = result | int(g.Taku)

	return newXmlElement(
		"GO",
		"type", strconv.Itoa(result),
		"lobby", g.Lobby,
	), nil
}

// "四南喰赤速" みたいな名前を出力
func (g GameInfo) Name() string {
	var name string
	if g.Sanma {
		name += "三"
	} else {
		name += "四"
	}
	name += g.Taku.ShortName()
	if g.Tonpu {
		name += "東"
	} else {
		name += "南"
	}
	if g.Ariari {
		name += "喰"
	}
	if g.Akadora {
		name += "赤"
	}
	if g.Sokutaku {
		name += "速"
	}
	return name
}

type Taku int

const (
	TakuIppan   Taku = 0x00
	TakuJoukyuu Taku = 0x80
	TakuTokujo  Taku = 0x20
	TakuHouou   Taku = 0xA0
)

func (r Taku) Name() string {
	switch r {
	case TakuIppan:
		return "一般"
	case TakuJoukyuu:
		return "上級"
	case TakuTokujo:
		return "特上"
	case TakuHouou:
		return "鳳凰"
	default:
		return "Unknown"
	}
}

func (r Taku) ShortName() string {
	switch r {
	case TakuIppan:
		return "般"
	case TakuJoukyuu:
		return "上"
	case TakuTokujo:
		return "特"
	case TakuHouou:
		return "鳳"
	default:
		return "Unknown"
	}
}

// Ensure GameInfo implements MJLogElement
var _ MJLogElement = &GameInfo{}
