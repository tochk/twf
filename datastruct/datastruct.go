package datastruct

type Field struct {
	Name              string
	Title             string
	Type              string
	IsNotDisabled     bool
	IsNotRequired     bool
	IsNotEditable     bool
	IsNotCreatable    bool
	IsNotShowOnList   bool
	IsNotShowOnItem   bool
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
