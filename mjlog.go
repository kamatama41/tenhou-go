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
