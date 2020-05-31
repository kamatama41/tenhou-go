package tenhou

type MJLog struct {
	Version  string
	Seed     string
	GameInfo GameInfo
	Players  Players
	Events   []MJLogElement
}

type MJLogElement interface {
	UnmarshalMJLog(e XmlElement) error
	MarshalMJLog() (XmlElement, error)
}

func (m *MJLog) GetResult() GameResult {
	// Eventsを逆に見ていってowariがある流局か和了を探す
	eventLen := len(m.Events)
	for i := eventLen - 1; i >= 0; i-- {
		switch e := m.Events[i].(type) {
		case *Agari:
			if o := e.Owari; !o.IsZero() {
				return o
			}
		case *Ryuukyoku:
			if o := e.Owari; !o.IsZero() {
				return o
			}
		}
	}
	return nil
}
