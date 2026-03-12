package valueobject

import "svc-boilerplate-fiber-golang/entity"

type Masterdata struct {
	entity.Masterdata
	entity.StandardKey
	entity.Pagination
	entity.Time
}

type MasterdataPayloadInsert struct {
	Data []Masterdata `json:"data" binding:"required"`
	User string
}

type MasterdataPayloadUpdate struct {
	Data []MasterdataDataUpdate `json:"data" binding:"required"`
	User string
}

type MasterdataDataUpdate struct {
	Param Masterdata `json:"param" binding:"required"`
	Body  Masterdata `json:"body" binding:"required"`
}

type MasterdataPayloadDelete struct {
	Param []Masterdata `json:"param" binding:"required"`
}
