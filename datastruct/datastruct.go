package datastruct

type Field struct {
	Name              string
	Title             string
	Type              string
	Placeholder       string
	Required          bool
	NoCreate          bool
	NoEdit            bool
	NotShowOnTable    bool
	ProcessParameters bool
	Value             string
	FkInfo            *FkInfo
	Disabled          bool
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
