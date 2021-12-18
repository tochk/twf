package datastruct

type Field struct {
	Name              string
	Title             string
	Type              string
	Placeholder       string
	IsNotDisabled     bool
	Required          bool
	NoCreate          bool
	NoEdit            bool
	NotShowOnTable    bool
	ProcessParameters bool
	Value             string
	FkInfo            *FkInfo
}

type FkInfo struct {
	FksIndex int
	ID       string
	Name     string
}

type FkKV struct {
	ID   interface{}
	Name interface{}
}
