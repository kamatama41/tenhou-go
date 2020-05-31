package tenhou

import "fmt"

type YakuType int

const (
	YakuTypeTsumohou       YakuType = 0
	YakuTypeRiichi         YakuType = 1
	YakuTypeIppatsu        YakuType = 2
	YakuTypeChankan        YakuType = 3
	YakuTypeRinshanKaihou  YakuType = 4
	YakuTypeHouteiRaoyui   YakuType = 5
	YakuTypeHaiteiRaoyui   YakuType = 6
	YakuTypePinfu          YakuType = 7
	YakuTypeTanyao         YakuType = 8
	YakuTypeIipeikou       YakuType = 9
	YakuTypeJikazeTon      YakuType = 10
	YakuTypeJikazeNan      YakuType = 11
	YakuTypeJikazeSha      YakuType = 12
	YakuTypeJikazePei      YakuType = 13
	YakuTypeBakazeTon      YakuType = 14
	YakuTypeBakazeNan      YakuType = 15
	YakuTypeBakazeSha      YakuType = 16
	YakuTypeBakazePei      YakuType = 17
	YakuTypeHaku           YakuType = 18
	YakuTypeHatsu          YakuType = 19
	YakuTypeChun           YakuType = 20
	YakuTypeDoubleReach    YakuType = 21
	YakuTypeChiitoitsu     YakuType = 22
	YakuTypeChanta         YakuType = 23
	YakuTypeIkkitsuukan    YakuType = 24
	YakuTypeSanshokuDoujun YakuType = 25
	YakuTypeSanshokuDoukou YakuType = 26
	YakuTypeSankantsu      YakuType = 27
	YakuTypeToitoi         YakuType = 28
	YakuTypeSanankou       YakuType = 29
	YakuTypeShousangen     YakuType = 30
	YakuTypeHonroutou      YakuType = 31
	YakuTypeRyanpeikou     YakuType = 32
	YakuTypeJunchan        YakuType = 33
	YakuTypeHonitsu        YakuType = 34
	YakuTypeChinitsu       YakuType = 35
	YakuTypeTenhou         YakuType = 37
	YakuTypeChiihou        YakuType = 38
	YakuTypeDaisangen      YakuType = 39
	YakuTypeSuuankou       YakuType = 40
	YakuTypeSuuankouTanki  YakuType = 41
	YakuTypeTsuuiisou      YakuType = 42
	YakuTypeRyuuiisou      YakuType = 43
	YakuTypeChinroutou     YakuType = 44
	YakuTypeChuurenPoutou  YakuType = 45
	YakuTypeChuurenPoutou9 YakuType = 46
	YakuTypeKokushiMusou   YakuType = 47
	YakuTypeKokushiMusou13 YakuType = 48
	YakuTypeDaisuushi      YakuType = 49
	YakuTypeShousuushi     YakuType = 50
	YakuTypeSuukantsu      YakuType = 51
	YakuTypeDora           YakuType = 52
	YakuTypeUradora        YakuType = 53
	YakuTypeAkadora        YakuType = 54
)

func NewYakuType(typ int) (YakuType, error) {
	for i := range yakuNames {
		if i == typ {
			return YakuType(typ), nil
		}
	}
	return 0, fmt.Errorf("unknown yaku type %d", typ)
}

var yakuNames = []string{
	"門前清自摸和", "立直", "一発", "槍槓", "嶺上開花",
	"海底摸月", "河底撈魚", "平和", "断幺九", "一盃口",
	"自風 東", "自風 南", "自風 西", "自風 北",
	"場風 東", "場風 南", "場風 西", "場風 北",
	"白", "發", "中",
	"両立直", "七対子", "混全帯幺九", "一気通貫",
	"三色同順", "三色同刻", "三槓子", "対々和", "三暗刻",
	"小三元", "混老頭", "二盃口", "純全帯幺九", "混一色",
	"清一色", "", "天和", "地和", "大三元",
	"四暗刻", "四暗刻単騎", "字一色", "緑一色", "清老頭",
	"九蓮宝燈", "純正九蓮宝燈", "国士無双", "国士無双１３面", "大四喜",
	"小四喜", "四槓子", "ドラ", "裏ドラ", "赤ドラ",
}

type Yaku struct {
	Type YakuType
	Han  int // 食い下がりや複数ドラを考慮
}

func (y Yaku) String() string {
	return fmt.Sprintf("%s(%d)", y.Name(), y.Han)
}

func (y Yaku) Name() string {
	return yakuNames[y.Type]
}

func (y Yaku) IsYakuman() bool {
	return y.Type >= YakuTypeTenhou && y.Type <= YakuTypeSuukantsu
}

// 赤, 表, 裏のいずれか
func (y Yaku) IsDora() bool {
	return y.Type >= YakuTypeDora || y.Type <= YakuTypeAkadora
}
