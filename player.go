package tenhou

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Players [4]Player

type Player struct {
	Name PlayerName
	Rate Rate
	Dan  Dan
	Sex  string
}

type PlayerIndex int

var AllPlayerIndexes = []PlayerIndex{0, 1, 2, 3}

func NewPlayerIndexFromString(pi string) (PlayerIndex, error) {
	idx, err := strconv.Atoi(pi)
	if err != nil {
		return 0, err
	}
	return NewPlayerIndex(idx)
}

func NewPlayerIndex(idx int) (PlayerIndex, error) {
	if idx < 0 || idx > 3 {
		return 0, fmt.Errorf("player index must be between 0 and 3 but %d", idx)
	}
	return PlayerIndex(idx), nil
}

func (p PlayerIndex) String() string {
	return strconv.Itoa(int(p))
}

func (p *Players) UnmarshalMJLog(e XmlElement) error {
	rates := splitByComma(e.AttrByName("rate"))
	dans, err := splitByCommaAsInt(e.AttrByName("dan"))
	if err != nil {
		return err
	}
	sexes := splitByComma(e.AttrByName("sx"))

	for _, pi := range AllPlayerIndexes {
		name, err := DecodePlayerName(e.AttrByName(fmt.Sprintf("n%d", pi)))
		if err != nil {
			return err
		}
		if name == "" {
			p[pi] = Player{
				Name: "",
				Rate: RateDefault,
				Dan:  0,
				Sex:  "C",
			}
			continue
		}
		rate := Rate{}
		if err := rate.Unmarshal(rates[pi]); err != nil {
			return err
		}
		dan := dans[pi]
		sex := sexes[pi]

		p[pi] = Player{
			Name: name,
			Rate: rate,
			Dan:  Dan(dan),
			Sex:  sex,
		}
	}
	return nil
}

func (p *Players) MarshalMJLog() (XmlElement, error) {
	elem := XmlElement{
		Name: "UN",
	}

	var dan []string
	var rate []string
	var sx []string
	for _, pi := range AllPlayerIndexes {
		// n0-n3
		player := p[pi]
		elem.AppendAttr(fmt.Sprintf("n%d", pi), player.Name.Encode())

		dan = append(dan, strconv.Itoa(int(player.Dan)))
		rate = append(rate, player.Rate.String())
		sx = append(sx, player.Sex)
	}
	// dan
	elem.AppendAttr("dan", joinByComma(dan))
	// rate
	elem.AppendAttr("rate", joinByComma(rate))
	// sx
	elem.AppendAttr("sx", joinByComma(sx))
	return elem, nil
}

type PlayerName string

func DecodePlayerName(encoded string) (PlayerName, error) {
	name, err := url.QueryUnescape(encoded)
	if err != nil {
		return "", err
	}
	return PlayerName(name), err
}

// https://gist.github.com/hnaohiro/4627658
func (n PlayerName) Encode() string {
	var result string
	for _, c := range n {
		if c <= 0x7f { // single byte
			result += fmt.Sprintf("%%%X", c)
		} else if c > 0x1fffff { // quaternary byte
			result += fmt.Sprintf("%%%X%%%X%%%X%%%X",
				0xf0+((c&0x1c0000)>>18),
				0x80+((c&0x3f000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else if c > 0x7ff { // triple byte
			result += fmt.Sprintf("%%%X%%%X%%%X",
				0xe0+((c&0xf000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else { // double byte
			result += fmt.Sprintf("%%%X%%%X",
				0xc0+((c&0x7c0)>>6),
				0x80+(c&0x3f),
			)
		}
	}
	return result
}

type Dan int

const (
	DanShinjin Dan = iota
	Dan9kyuu
	Dan8kyuu
	Dan7kyuu
	Dan6kyuu
	Dan5kyuu
	Dan4kyuu
	Dan3kyuu
	Dan2kyuu
	Dan1kyuu
	Dan1dan
	Dan2dan
	Dan3dan
	Dan4dan
	Dan5dan
	Dan6dan
	Dan7dan
	Dan8dan
	Dan9dan
	Dan10dan
	DanTenhou
)

func (d Dan) String() string {
	return strconv.Itoa(int(d))
}

func (d Dan) Name() string {
	switch d {
	case DanShinjin:
		return "新人"
	case Dan9kyuu:
		return "９級"
	case Dan8kyuu:
		return "８級"
	case Dan7kyuu:
		return "７級"
	case Dan6kyuu:
		return "６級"
	case Dan5kyuu:
		return "５級"
	case Dan4kyuu:
		return "４級"
	case Dan3kyuu:
		return "３級"
	case Dan2kyuu:
		return "２級"
	case Dan1kyuu:
		return "１級"
	case Dan1dan:
		return "初段"
	case Dan2dan:
		return "二段"
	case Dan3dan:
		return "三段"
	case Dan4dan:
		return "四段"
	case Dan5dan:
		return "五段"
	case Dan6dan:
		return "六段"
	case Dan7dan:
		return "七段"
	case Dan8dan:
		return "八段"
	case Dan9dan:
		return "九段"
	case Dan10dan:
		return "十段"
	case DanTenhou:
		return "天鳳"
	default:
		return "Unknown"
	}
}

// rateは"1947.54" のように小数点付きで来るので小数も一応保持しておく
type Rate struct {
	I int
	D int
}

var RateDefault = Rate{1500, 0}

func (r *Rate) String() string {
	return fmt.Sprintf("%d.%02d", r.I, r.D)
}

func (r *Rate) Unmarshal(rate string) error {
	split := strings.Split(rate, ".")
	i, err := strconv.Atoi(split[0])
	if err != nil {
		return err
	}
	r.I = i
	d, err := strconv.Atoi((split[1] + "0")[:2])
	if err != nil {
		return err
	}
	r.D = d
	return nil
}

// Ensure Players implements MJLogElement
var _ MJLogElement = &Players{}
