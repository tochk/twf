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

type Config struct {
	Title       string
	IsEditable  bool
	IsDeletable bool
	IsCreatable bool
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
