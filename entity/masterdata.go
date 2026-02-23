package entity

type MasterdataKey struct {
	ForeignKey uint64 `json:"-"`
}

type Masterdata struct {
	MasterdataKey
	Column string `json:"column,omitempty"`
}
