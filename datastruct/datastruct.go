// Package datastruct holds the field descriptions twf builds from struct tags.
package datastruct

// Field describes one struct field parsed from its twf tag.
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

// FkInfo holds the parsed fk:<index>;<idField>;<nameField> tag option.
type FkInfo struct {
	FksIndex int
	ID       string
	Name     string
}

// FkKV is a resolved foreign key: the id value and its display name.
type FkKV struct {
	ID   any
	Name any
}
