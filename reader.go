package tenhou

import (
	"fmt"
	"strings"
)

type Reader struct {
	ml *MJLog
	p  Printer
}

type Printer interface {
	Printf(msg string, a ...interface{})
}

func NewReader(mjlog *MJLog, printer Printer) *Reader {
	return &Reader{
		ml: mjlog,
		p:  printer,
	}
}

func (r *Reader) ReadAll() {
	r.p.Printf("ルール: %s", r.ml.GameInfo.Name())
	r.p.Printf("プレイヤー")
	for _, p := range r.ml.Players {
		r.p.Printf(" %s %s %d %s", p.Name, p.Dan.Name(), p.Rate.I, p.Sex)
	}
	r.p.Printf("対局開始")
	for _, event := range r.ml.Events {
		r.printEvent(event)
	}
}

func (r *Reader) printEvent(e MJLogElement) {
	switch e := e.(type) {
	case *KyokuStart:
		r.p.Printf("========================================")
		r.p.Printf("[%s] %d本場 供託%d本 ドラ表示牌%s 親 %s", e.Round.Name(), e.Honba, e.Kyotaku, e.Dora, r.pn(e.Oya))
		r.p.Printf("配牌")
		for i, haiPai := range e.HaiPai {
			r.p.Printf(" %s %s", r.pn(i), haiPai.Sort())
		}
	case *Tsumo:
		r.p.Printf("%sが%sをツモ", r.pn(e.Who), e.Hai)
	case *Dahai:
		r.p.Printf("%sが%sを打牌", r.pn(e.Who), e.Hai)
	case *Naki:
		r.printNaki(e)
	case *DoraOpen:
		r.p.Printf("新ドラ表示 %s", e.Hai.Name())
	case *Reach:
		if e.Step == 1 {
			r.p.Printf("%sが立直宣言", r.pn(e.Who))
		} else {
			r.p.Printf("%sのリーチが成立。リーチ後の得点: %s", r.pn(e.Who), e.Ten)
		}
	case *Bye:
		r.p.Printf("%sが接続断", r.pn(e.Who))
	case *Reconnect:
		r.p.Printf("%sが復帰", e.Name)
	case *Agari:
		r.printAgari(e)
	case *Ryuukyoku:
		r.printRyuKyoku(e)
	}
}

func (r *Reader) printNaki(naki *Naki) {
	msg := fmt.Sprintf("%sが", r.pn(naki.Who))
	mentsu := naki.Mentsu
	from := mentsu.From().FromWho(naki.Who)
	// 自分じゃないときは相手を表示
	if naki.Who != from {
		msg += fmt.Sprintf("%sから", r.pn(from))
	}
	msg += fmt.Sprintf("%s %s", mentsu.Name(), mentsu.Mentsu())

	r.p.Printf(msg)
}

func (r *Reader) printAgari(agari *Agari) {
	msg := fmt.Sprintf("%sが", r.pn(agari.Who))
	if agari.IsRon() {
		msg += "ロン和了"
	} else {
		msg += "ツモ和了"
	}
	r.p.Printf(msg)

	msg = fmt.Sprintf(" 手牌: %s", agari.Tehai)
	if len(agari.NakiMentsu) > 0 {
		var fmStr []string
		for _, m := range agari.NakiMentsu {
			fmStr = append(fmStr, fmt.Sprintf("[%s]", m.Mentsu().String()))
		}
		msg += fmt.Sprintf(" %s", strings.Join(fmStr, " "))
	}
	r.p.Printf(msg)

	msg = " ドラ表示牌:"
	for _, ud := range agari.Dora {
		msg += " " + ud.Name()
	}
	r.p.Printf(msg)

	if len(agari.Uradora) > 0 {
		msg = " 裏ドラ表示牌:"
		for _, ud := range agari.Uradora {
			msg += " " + ud.Name()
		}
		r.p.Printf(msg)
	}

	msg = " 和了役:"
	for _, y := range agari.Yaku {
		msg += " " + y.Name()
		if y.IsDora() {
			msg += fmt.Sprintf("%d", y.Han)
		}
	}
	r.p.Printf(msg)

	r.p.Printf(" 得点: %d翻%d符 %d点", agari.Han, agari.Fu, agari.Ten)
	r.p.Printf(" 積み棒%d本 リーチ棒%d本", agari.TsumiBo, agari.ReachBo)
	r.p.Printf(" 点数移動: %s", agari.SC)
	if agari.Owari != "" {
		r.p.Printf("終局 %s", agari.Owari)
	}
}

func (r *Reader) printRyuKyoku(ryuKyoku *Ryuukyoku) {
	r.p.Printf("流局")
	if ryuKyoku.Type != "" {
		r.p.Printf(" 理由: %s", ryuKyoku.Type.Name())
	}
	r.p.Printf(" 点数移動: %s", ryuKyoku.SC)
	if ryuKyoku.Owari != "" {
		r.p.Printf("終局 %s", ryuKyoku.Owari)
	}
}

// PlayerNameのエイリアス
func (r *Reader) pn(i PlayerIndex) PlayerName {
	return r.ml.Players[i].Name
}
