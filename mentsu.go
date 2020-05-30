package tenhou

import (
	"fmt"
)

type NakiFrom int

const (
	NakiFromSelf NakiFrom = iota // 暗槓と抜きドラのとき
	NakiFromShimocha
	NakiFromToimen
	NakiFromKamicha
)

// whoから見て誰がやったか
func (f NakiFrom) FromWho(who PlayerIndex) PlayerIndex {
	switch f {
	case NakiFromSelf:
		return who
	case NakiFromShimocha:
		return (who + 1) % 4
	case NakiFromToimen:
		return (who + 2) % 4
	case NakiFromKamicha:
		return (who + 3) % 4
	default:
		return -1
	}
}

func (f NakiFrom) String() string {
	switch f {
	case NakiFromSelf:
		return ""
	case NakiFromShimocha:
		return "下家"
	case NakiFromToimen:
		return "対面"
	case NakiFromKamicha:
		return "上家"
	default:
		return "unknown"
	}
}

type NakiMentsu interface {
	Name() string
	From() NakiFrom
	Mentsu() HaiList
	What() Hai
}

type Chii struct {
	mentsu HaiList
	what   int // 順子の中のどれを鳴いたか
}

func (c Chii) Name() string {
	return "チー"
}

func (c Chii) From() NakiFrom {
	return NakiFromKamicha
}

func (c Chii) Mentsu() HaiList {
	return c.mentsu
}

func (c Chii) What() Hai {
	return c.mentsu[c.what]
}

type Pon struct {
	from   NakiFrom
	mentsu HaiList
	what   int // 刻子の中のどれを鳴いたか
}

func (p Pon) Name() string {
	return "ポン"
}

func (p Pon) From() NakiFrom {
	return p.from
}

func (p Pon) Mentsu() HaiList {
	return p.mentsu
}

func (p Pon) What() Hai {
	return p.mentsu[p.what]
}

type Ankan struct {
	what Hai
}

func (ak Ankan) Name() string {
	return "暗槓"
}

func (ak Ankan) From() NakiFrom {
	return NakiFromSelf
}

func (ak Ankan) Mentsu() HaiList {
	var res HaiList
	for i := 0; i < 4; i++ {
		res = append(res, Hai{
			Type:   ak.what.Type,
			Number: ak.what.Number,
			Index:  i,
		})
	}
	return res
}

func (ak Ankan) What() Hai {
	return ak.what
}

type Minkan struct {
	from NakiFrom
	what Hai
}

func (mk Minkan) Name() string {
	return "明槓"
}

func (mk Minkan) From() NakiFrom {
	return mk.from
}

func (mk Minkan) Mentsu() HaiList {
	var res HaiList
	for i := 0; i < 4; i++ {
		res = append(res, Hai{
			Type:   mk.what.Type,
			Number: mk.what.Number,
			Index:  i,
		})
	}
	return res
}

func (mk Minkan) What() Hai {
	return mk.what
}

type Chakan struct {
	pon  Pon // 槓する前のポン情報
	what Hai
}

func (ck Chakan) Name() string {
	return "加槓"
}

func (ck Chakan) From() NakiFrom {
	return NakiFromSelf
}

func (ck Chakan) Mentsu() HaiList {
	var res HaiList
	for i := 0; i < 4; i++ {
		res = append(res, Hai{
			Type:   ck.what.Type,
			Number: ck.what.Number,
			Index:  i,
		})
	}
	return res
}

func (ck Chakan) What() Hai {
	return ck.what
}

type Nukidora struct {
	what Hai
}

func (nd Nukidora) Name() string {
	return "抜きドラ"
}

func (nd Nukidora) From() NakiFrom {
	return NakiFromSelf
}

func (nd Nukidora) Mentsu() HaiList {
	return HaiList{nd.what}
}

func (nd Nukidora) What() Hai {
	return nd.what
}

/**
 * 副露面子
 * ref: https://blog.kobalab.net/entry/20170228/1488294993#f-3877133a
 * ref: http://tenhou.net/img/tehai.js
 */
func UnmarshalMentsu(m int) (NakiMentsu, error) {
	// 0x0003 (0b0000000000000011): 誰から鳴いたか
	from := NakiFrom(m & 0x0003)
	// 0x0004 (0b0000000000000100): 面子タイプ
	mentsuType := (m & 0x0004) >> 2

	if mentsuType == 1 {
		// 順子 (チー)
		// 0xFC00 (0b1111110000000000): 面子タイプ
		patternInt := (m & 0xFC00) >> 10
		// 何番目の牌を鳴いたか
		what := patternInt % 3
		patternInt = patternInt / 3
		// 萬子、筒子、索子を決定
		haiType := patternInt / 7
		// 形を決定 (0:123 ~ 6:789)
		startNum := patternInt % 7
		// 牌添字1~3を取得
		haiIndexes := []int{
			(m & 0x0018) >> 3, // 0b0000000000011000
			(m & 0x0060) >> 5, // 0b0000000001100000
			(m & 0x0180) >> 7, // 0b0000000110000000
		}

		chi := Chii{what: what}
		// 面子に変換
		for i, idx := range haiIndexes {
			hai := Hai{
				Type:   HaiType(haiType),
				Number: startNum + i,
				Index:  idx,
			}
			chi.mentsu = append(chi.mentsu, hai)
		}
		return chi, nil
	} else {
		// 刻子 (ポン or カン) 0b0000000000011000
		koutsuType := (m & 0x0018) >> 3
		if koutsuType == 0 {
			// 暗槓 or 大明槓
			// 0xFF00 (0b1111111100000000): 槓子タイプ
			patternInt := (m & 0xFF00) >> 8
			// どの牌を鳴いたか
			// (暗槓の場合はこの数字は意味無さそう)
			idxWhat := patternInt % 4
			patternInt = patternInt / 4
			// 萬子、筒子、索子、字牌を決定
			haiType := patternInt / 9
			// 数字を決定 (字牌は0-6)
			haiNumber := patternInt % 9
			// 6ビット目: 抜きドラかどうか (抜きドラのとき1)
			isNukiDora := (m&0x0020)>>5 == 1

			hai := Hai{
				Type:   HaiType(haiType),
				Number: haiNumber,
				Index:  idxWhat,
			}
			if isNukiDora {
				return Nukidora{hai}, nil
			} else if from == NakiFromSelf {
				return Ankan{hai}, nil
			} else {
				return Minkan{from, hai}, nil
			}
		} else if koutsuType == 1 || koutsuType == 2 {
			// 0xFE00 (0b1111111000000000): 刻子タイプ
			patternInt := (m & 0xFE00) >> 9
			// 何番目の牌を鳴いたか (最終形3枚を昇順に並べた場合のindex)
			what := patternInt % 3
			patternInt = patternInt / 3
			// 萬子、筒子、索子、字牌を決定
			haiType := patternInt / 9
			// 数字を決定 (字牌は0-6)
			haiNumber := patternInt % 9
			// ポンで鳴かなかった牌のindex (0b0000000001100000)
			haiIndex := m & 0x0060 >> 5

			pon := Pon{from: from, what: what}
			for i := 0; i < 4; i++ {
				// 牌添字に指定されている牌はポンされていないので省く
				if i == haiIndex {
					continue
				} else {
					hai := Hai{
						Type:   HaiType(haiType),
						Number: haiNumber,
						Index:  i,
					}
					pon.mentsu = append(pon.mentsu, hai)
				}
			}

			if koutsuType == 1 {
				// ポン
				return pon, nil
			} else {
				hai := Hai{
					Type:   HaiType(haiType),
					Number: haiNumber,
					Index:  haiIndex,
				}
				// 加槓
				return Chakan{pon, hai}, nil
			}
		} else {
			return nil, fmt.Errorf("unexpected Kotsu type %d m:%d", koutsuType, m)
		}
	}
}

func MarshalMentsu(men NakiMentsu) int {
	var m int
	switch men := men.(type) {
	case Chii:
		// 必ず上家から鳴く
		m = m | int(NakiFromKamicha)
		// 面子タイプ=1
		m = m | 1<<2
		// 肺添字を取得
		for i, hai := range men.mentsu {
			m = m | hai.Index<<((2*i)+3)
		}
		// 牌パターンをエンコード
		haiStart := men.mentsu[0]
		haiType := int(haiStart.Type)
		startNum := haiStart.Number
		patternInt := (haiType*7+startNum)*3 + men.what
		m = m | patternInt<<10
	case Ankan, Minkan, Nukidora:
		// 誰から鳴いたか (暗槓の場合は自分)
		m = m | int(men.From())
		// 牌パターンをエンコード
		haiType := int(men.What().Type)
		number := men.What().Number
		patternInt := (haiType*9+number)*4 + men.What().Index
		m = m | patternInt<<8
		// 抜きドラの場合はフラグを立てる
		switch men.(type) {
		case Nukidora:
			m = m | 1<<5
		}
	case Pon:
		// 誰から鳴いたか
		m = m | int(men.from)
		// 刻子タイプ=1
		m = m | 1<<3
		// ポンしなかった牌のindexを探す
		var idxRest int
	LOOP:
		for i := 0; i < 4; i++ {
			for _, h := range men.mentsu {
				if h.Index == i {
					continue LOOP
				}
			}
			idxRest = i
			break
		}
		m = m | idxRest<<5
		// 牌パターンをエンコード
		what := men.mentsu[men.what]
		haiType := int(what.Type)
		number := what.Number
		patternInt := (haiType*9+number)*3 + men.what
		m = m | patternInt<<9
	case Chakan:
		// ポンを誰から鳴いたか
		m = m | int(men.pon.from)
		// 刻子タイプ=2
		m = m | 2<<3
		// 槓した牌のindex
		m = m | men.what.Index<<5

		// 牌パターンをエンコード
		haiType := int(men.what.Type)
		number := men.what.Number
		patternInt := (haiType*9+number)*3 + men.pon.what
		m = m | patternInt<<9
	}
	return m
}
