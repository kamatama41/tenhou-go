package tenhou

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// 局開始
type KyokuStart struct {
	Round   Round
	Honba   int
	Kyotaku int
	Dora    Hai
	Ten     []int
	Oya     PlayerIndex
	HaiPai  map[PlayerIndex]HaiList
	Dice    [2]int
}

type Round int

func (r Round) Name() string {
	var wind string
	switch r / 4 {
	case 0:
		wind = "東"
	case 1:
		wind = "南"
	case 2:
		wind = "西"
	case 3:
		wind = "北" // 存在する?
	}
	return fmt.Sprintf("%s%d局", wind, r%4+1)
}

func (k *KyokuStart) UnmarshalMJLog(e XmlElement) error {
	return e.ForEachAttr(func(name, value string) error {
		switch name {
		case "seed":
			seed, err := splitByCommaAsInt(value)
			if err != nil {
				return err
			}
			k.Round = Round(seed[0])
			k.Honba = seed[1]
			k.Kyotaku = seed[2]
			// 3,4は開門のサイコロの値 (どこで使う?)
			k.Dice[0] = seed[3]
			k.Dice[1] = seed[4]
			if err := k.Dora.Unmarshal(seed[5]); err != nil {
				return err
			}
		case "ten":
			ten, err := splitByCommaAsInt(value)
			if err != nil {
				return err
			}
			k.Ten = ten
		case "oya":
			oya, err := NewPlayerIndexFromString(value)
			if err != nil {
				return err
			}
			k.Oya = oya
		case "hai0", "hai1", "hai2", "hai3":
			if value == "" {
				// 三麻のときはhai3のvalueが空
				return nil
			}
			who, err := NewPlayerIndexFromString(name[3:])
			if err != nil {
				return err
			}
			tehai, err := splitByCommaAsHaiList(value)
			if err != nil {
				return err
			}
			if k.HaiPai == nil {
				k.HaiPai = make(map[PlayerIndex]HaiList)
			}
			k.HaiPai[who] = tehai
		default:
			return newUnknownAttrErr(name, value)
		}
		return nil
	})
}

func (k *KyokuStart) MarshalMJLog() (XmlElement, error) {
	elem := newXmlElement("INIT")
	// seed
	seed := []int{
		int(k.Round),
		k.Honba,
		k.Kyotaku,
		k.Dice[0],
		k.Dice[1],
		k.Dora.ID(),
	}
	elem.AppendAttr("seed", joinByCommaFromInts(seed))
	elem.AppendAttr("ten", joinByCommaFromInts(k.Ten))
	elem.AppendAttr("oya", k.Oya.String())
	for _, pi := range AllPlayerIndexes {
		elem.AppendAttr(fmt.Sprintf("hai%d", pi), joinByCommaFromInts(k.HaiPai[pi].IDs()))
	}

	return elem, nil
}

// ツモ
type Tsumo struct {
	Who PlayerIndex
	Hai Hai
}

func (t *Tsumo) UnmarshalMJLog(e XmlElement) error {
	who := strings.Index("TUVW", e.Name[0:1])
	if who == -1 {
		return fmt.Errorf("index of %s not found", e.Name)
	}
	t.Who = PlayerIndex(who)
	hai, err := strconv.Atoi(e.Name[1:])
	if err != nil {
		return err
	}
	if err := t.Hai.Unmarshal(hai); err != nil {
		return err
	}
	return nil
}

func (t *Tsumo) MarshalMJLog() (XmlElement, error) {
	i := int(t.Who)
	name := fmt.Sprintf("%s%d", "TUVW"[i:i+1], t.Hai.ID())
	return newXmlElement(name), nil
}

// 打牌
type Dahai struct {
	Who PlayerIndex
	Hai Hai
}

func (d *Dahai) UnmarshalMJLog(e XmlElement) error {
	who, err := NewPlayerIndex(strings.Index("DEFG", e.Name[0:1]))
	if err != nil {
		return err
	}
	d.Who = who

	hai, err := strconv.Atoi(e.Name[1:])
	if err != nil {
		return err
	}
	if err := d.Hai.Unmarshal(hai); err != nil {
		return err
	}

	return nil
}

func (d *Dahai) MarshalMJLog() (XmlElement, error) {
	i := int(d.Who)
	name := fmt.Sprintf("%s%d", "DEFG"[i:i+1], d.Hai.ID())
	return newXmlElement(name), nil
}

// 鳴き
type Naki struct {
	Who    PlayerIndex
	Mentsu NakiMentsu
}

func (n *Naki) UnmarshalMJLog(e XmlElement) error {
	return e.ForEachAttr(func(name, value string) error {
		switch name {
		case "who":
			who, err := NewPlayerIndexFromString(value)
			if err != nil {
				return err
			}
			n.Who = who
		case "m":
			m, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			mentsu, err := UnmarshalMentsu(m)
			if err != nil {
				return err
			}
			n.Mentsu = mentsu
		default:
			return newUnknownAttrErr(name, value)
		}
		return nil
	})
}

func (n *Naki) MarshalMJLog() (XmlElement, error) {
	elem := newXmlElement("N")
	elem.AppendAttr("who", n.Who.String())
	elem.AppendAttr("m", strconv.Itoa(MarshalMentsu(n.Mentsu)))
	return elem, nil
}

// ドラ表示牌をめくる
type DoraOpen struct {
	Hai Hai
}

func (d *DoraOpen) UnmarshalMJLog(e XmlElement) error {
	return e.ForEachAttr(func(name, value string) error {
		switch name {
		case "hai":
			haiId, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			if err := d.Hai.Unmarshal(haiId); err != nil {
				return err
			}
		default:
			return newUnknownAttrErr(name, value)
		}
		return nil
	})
}

func (d *DoraOpen) MarshalMJLog() (XmlElement, error) {
	elem := newXmlElement("DORA")
	elem.AppendAttr("hai", strconv.Itoa(d.Hai.ID()))
	return elem, nil
}

type Reach struct {
	Who  PlayerIndex
	Step int // 1だと立直宣言, 2だと立直成立
	Ten  string
}

func (r *Reach) UnmarshalMJLog(e XmlElement) error {
	return e.ForEachAttr(func(name, value string) error {
		switch name {
		case "who":
			who, err := NewPlayerIndexFromString(value)
			if err != nil {
				return err
			}
			r.Who = who
		case "step":
			step, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			r.Step = step
		case "ten":
			// 立直後の全員の得点, カンマ区切り
			r.Ten = value
		default:
			return newUnknownAttrErr(name, value)
		}
		return nil
	})
}

func (r *Reach) MarshalMJLog() (XmlElement, error) {
	elem := newXmlElement("REACH")
	elem.AppendAttr("who", r.Who.String())
	if r.Ten != "" {
		elem.AppendAttr("ten", r.Ten)
	}
	elem.AppendAttr("step", strconv.Itoa(r.Step))
	return elem, nil
}

// 接続断
type Bye struct {
	Who PlayerIndex
}

func (b *Bye) UnmarshalMJLog(e XmlElement) error {
	return e.ForEachAttr(func(name, value string) error {
		switch name {
		case "who":
			who, err := NewPlayerIndexFromString(value)
			if err != nil {
				return err
			}
			b.Who = who
		default:
			return newUnknownAttrErr(name, value)
		}
		return nil
	})
}

func (b *Bye) MarshalMJLog() (XmlElement, error) {
	elem := newXmlElement("BYE")
	elem.AppendAttr("who", b.Who.String())
	return elem, nil
}

// 再接続
type Reconnect struct {
	Who  PlayerIndex
	Name PlayerName
}

func (r *Reconnect) UnmarshalMJLog(e XmlElement) error {
	return e.ForEachAttr(func(name, value string) error {
		switch name {
		case "n0", "n1", "n2", "n3":
			if r.Name != "" {
				return fmt.Errorf("reconnected user already registered")
			}
			who, err := NewPlayerIndexFromString(name[1:])
			if err != nil {
				return err
			}
			name, err := DecodePlayerName(value)
			if err != nil {
				return err
			}
			r.Who = who
			r.Name = name
		default:
			return newUnknownAttrErr(name, value)
		}
		return nil
	})
}

func (r *Reconnect) MarshalMJLog() (XmlElement, error) {
	elem := newXmlElement("UN")
	elem.AppendAttr(fmt.Sprintf("n%s", r.Who), r.Name.Encode())
	return elem, nil
}

// 和了
type Agari struct {
	Who        PlayerIndex
	From       PlayerIndex
	Han        int
	Yaku       []Yaku
	Tehai      HaiList
	NakiMentsu []NakiMentsu
	Agari      Hai
	Dora       HaiList
	Uradora    HaiList
	Fu         int
	Ten        int
	ManganType int
	TsumiBo    int
	ReachBo    int
	SC         string
	Owari      GameResult
}

func (a *Agari) IsRon() bool {
	return a.Who != a.From
}

func (a *Agari) UnmarshalMJLog(e XmlElement) error {
	return e.ForEachAttr(func(name, value string) error {
		switch name {
		case "ba":
			// 積み棒とリーチ棒
			tsumiBoAndReachBo, err := splitByCommaAsInt(value)
			if err != nil {
				return err
			}
			a.TsumiBo = tsumiBoAndReachBo[0]
			a.ReachBo = tsumiBoAndReachBo[1]
		case "hai":
			// 手牌情報
			teHai, err := splitByCommaAsHaiList(value)
			if err != nil {
				return err
			}
			a.Tehai = teHai
		case "m":
			// 副露面子
			ms, err := splitByCommaAsInt(value)
			if err != nil {
				return err
			}
			for _, m := range ms {
				mentsu, err := UnmarshalMentsu(m)
				if err != nil {
					return err
				}
				a.NakiMentsu = append(a.NakiMentsu, mentsu)
			}
		case "machi":
			// 和了牌
			agari, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			if err := a.Agari.Unmarshal(agari); err != nil {
				return err
			}
		case "ten":
			// 0:符, 1:和了点, 2:満貫情報 (0: 満貫未満、1: 満貫、2: 跳満、3: 倍満、4: 三倍満、5: 役満)
			fuAndTenAndManganType, err := splitByCommaAsInt(value)
			if err != nil {
				return err
			}
			a.Fu = fuAndTenAndManganType[0]
			a.Ten = fuAndTenAndManganType[1]
			a.ManganType = fuAndTenAndManganType[2]
		case "yaku":
			// "役Type, 飜数, 役Type,　飜数..." で並んでいる
			yakuTypeAndHans, err := splitByCommaAsInt(value)
			hanSum := 0
			if err != nil {
				return err
			}
			for i := 0; i < len(yakuTypeAndHans); i += 2 {
				id := yakuTypeAndHans[i]
				han := yakuTypeAndHans[i+1]
				yakuType, err := NewYakuType(id)
				if err != nil {
					return err
				}
				a.Yaku = append(a.Yaku, Yaku{yakuType, han})
				hanSum += han
			}
			a.Han = hanSum
		case "doraHai":
			// ドラ
			dora, err := splitByCommaAsHaiList(value)
			if err != nil {
				return err
			}
			a.Dora = dora
		case "who":
			// 誰があがったか
			who, err := NewPlayerIndexFromString(value)
			if err != nil {
				return err
			}
			a.Who = who
		case "fromWho":
			// 誰からあがったか (ツモの場合は自分)
			from, err := NewPlayerIndexFromString(value)
			if err != nil {
				return err
			}
			a.From = from
		case "yakuman":
			// 役満はIDのみ
			yakuTypes, err := splitByCommaAsInt(value)
			if err != nil {
				return err
			}
			for _, typ := range yakuTypes {
				yakuType, err := NewYakuType(typ)
				if err != nil {
					return err
				}
				a.Yaku = append(a.Yaku, Yaku{yakuType, 13})
			}
		case "doraHaiUra":
			// 裏ドラ
			uraDora, err := splitByCommaAsHaiList(value)
			if err != nil {
				return err
			}
			a.Uradora = uraDora
		case "sc":
			// 局収支
			//  "0の和了前得点,0の今回の得失点,1の和了前得点,..."
			// 使わないのでとりあえず保存だけ
			a.SC = value
		case "owari":
			// 対局が終わったときに入る最終的な得点とポイント
			if err := a.Owari.Unmarshal(value); err != nil {
				return err
			}
		default:
			return newUnknownAttrErr(name, value)
		}

		return nil
	})
}

func (a *Agari) MarshalMJLog() (XmlElement, error) {
	elem := newXmlElement("AGARI")
	ba := []int{
		a.TsumiBo,
		a.ReachBo,
	}
	elem.AppendAttr("ba", joinByCommaFromInts(ba))
	elem.AppendAttr("hai", joinByCommaFromInts(a.Tehai.IDs()))
	if len(a.NakiMentsu) > 0 {
		var m []int
		for _, men := range a.NakiMentsu {
			m = append(m, MarshalMentsu(men))
		}
		elem.AppendAttr("m", joinByCommaFromInts(m))
	}
	elem.AppendAttr("machi", strconv.Itoa(a.Agari.ID()))
	// ten
	ten := []int{
		a.Fu,
		a.Ten,
		a.ManganType,
	}
	elem.AppendAttr("ten", joinByCommaFromInts(ten))
	// yaku
	var yaku []int
	for _, y := range a.Yaku {
		yaku = append(yaku, int(y.Type), y.Han)
	}
	elem.AppendAttr("yaku", joinByCommaFromInts(yaku))
	elem.AppendAttr("doraHai", joinByCommaFromInts(a.Dora.IDs()))
	if len(a.Uradora) > 0 {
		elem.AppendAttr("doraHaiUra", joinByCommaFromInts(a.Uradora.IDs()))
	}
	elem.AppendAttr("who", a.Who.String())
	elem.AppendAttr("fromWho", a.From.String())
	elem.AppendAttr("sc", a.SC)
	if !a.Owari.IsZero() {
		elem.AppendAttr("owari", a.Owari.Marshal())
	}
	return elem, nil
}

// 流局
type Ryuukyoku struct {
	TsumiBo int
	ReachBo int
	Tenpai  map[PlayerIndex]HaiList
	Type    RyuukyokuType
	SC      string
	Owari   GameResult
}

type RyuukyokuType string

const (
	RyuukyokuTypeNM     RyuukyokuType = "nm"
	RyuukyokuTypeYao9   RyuukyokuType = "yao9"
	RyuukyokuTypeKaze4  RyuukyokuType = "kaze4"
	RyuukyokuTypeReach4 RyuukyokuType = "reach4"
	RyuukyokuTypeRon3   RyuukyokuType = "ron3"
	RyuukyokuTypeKan4   RyuukyokuType = "kan4"
)

func (rt RyuukyokuType) Name() string {
	switch rt {
	case RyuukyokuTypeNM:
		return "流し満貫"
	case RyuukyokuTypeYao9:
		return "九種九牌"
	case RyuukyokuTypeKaze4:
		return "四風連打"
	case RyuukyokuTypeReach4:
		return "四家立直"
	case RyuukyokuTypeRon3:
		return "三家和了"
	case RyuukyokuTypeKan4:
		return "四槓散了"
	default:
		return string(rt)
	}
}

func (r *Ryuukyoku) UnmarshalMJLog(e XmlElement) error {
	return e.ForEachAttr(func(name, value string) error {
		switch name {
		case "type":
			r.Type = RyuukyokuType(value)
		case "ba":
			// 積み棒、リーチ棒
			tsumiboAndReachBo, err := splitByCommaAsInt(value)
			if err != nil {
				return err
			}
			r.TsumiBo = tsumiboAndReachBo[0]
			r.ReachBo = tsumiboAndReachBo[1]
		case "hai0", "hai1", "hai2", "hai3":
			// 聴牌している人の分だけ入っている
			who, err := NewPlayerIndexFromString(name[3:])
			if err != nil {
				return err
			}
			tehai, err := splitByCommaAsHaiList(value)
			if err != nil {
				return err
			}
			if r.Tenpai == nil {
				r.Tenpai = make(map[PlayerIndex]HaiList)
			}
			r.Tenpai[who] = tehai
		case "sc":
			// agariのscと同じ
			r.SC = value
		case "owari":
			// agariのowariと同じ
			if err := r.Owari.Unmarshal(value); err != nil {
				return err
			}
		default:
			return newUnknownAttrErr(name, value)
		}
		return nil
	})
}

func (r *Ryuukyoku) MarshalMJLog() (XmlElement, error) {
	elem := newXmlElement("RYUUKYOKU")
	ba := []int{
		r.TsumiBo,
		r.ReachBo,
	}
	if r.Type != "" {
		elem.AppendAttr("type", string(r.Type))
	}
	elem.AppendAttr("ba", joinByCommaFromInts(ba))
	elem.AppendAttr("sc", r.SC)
	for _, pi := range AllPlayerIndexes {
		hai, ok := r.Tenpai[pi]
		if !ok {
			continue
		}
		elem.AppendAttr(fmt.Sprintf("hai%d", pi), joinByCommaFromInts(hai.IDs()))
	}
	if !r.Owari.IsZero() {
		elem.AppendAttr("owari", r.Owari.Marshal())
	}
	return elem, nil
}

type GameResult []struct {
	Player PlayerIndex
	Ten    int
	Point  Point
}

// "0の得点,0のポイント,1の得点,..."の形式 (pointは100で割った値が入っている)
func (g *GameResult) Unmarshal(owari string) error {
	tenAndPoint := splitByComma(owari)
	for i := 0; i < len(tenAndPoint); i += 2 {
		ten, err := strconv.Atoi(tenAndPoint[i])
		if err != nil {
			return err
		}
		point, err := strconv.ParseFloat(tenAndPoint[i+1], 64)
		if err != nil {
			return err
		}
		pi, err := NewPlayerIndex(i / 2)
		if err != nil {
			return err
		}
		// サンマの4人目は "0,0" になるため除外
		if ten == 0 && point == 0 {
			continue
		}
		*g = append(*g, struct {
			Player PlayerIndex
			Ten    int
			Point  Point
		}{pi, ten * 100, Point(point)})
	}
	return nil
}

func (g GameResult) Marshal() string {
	var res []string
	for _, r := range g {
		res = append(res, fmt.Sprintf("%d,%s", r.Ten/100, r.Point.String()))
	}
	// サンマのときは4人目を追加
	if len(g) == 3 {
		res = append(res, "0,0")
	}
	return strings.Join(res, ",")
}

func (g GameResult) IsZero() bool {
	return len(g) == 0
}

// ポイント順でソートする
func (g GameResult) Sort() GameResult {
	var res GameResult
	res = append(res, g...)
	sort.Slice(res, func(i, j int) bool { return res[i].Point > res[j].Point })
	return res
}

type Point float64

func (p Point) String() string {
	return fmt.Sprintf("%.1f", p)
}
