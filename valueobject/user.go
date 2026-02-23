package valueobject

import "boilerplate-golang/entity"

type User struct {
	entity.User
	entity.Identity
	entity.StandardKey
	entity.Pagination
	entity.Time
}

type UserPayloadInsert struct {
	Data User `json:"users" binding:"required"`
	User string
}

type UserPayloadUpdateArray struct {
	Data []UsersDataUpdate `json:"users" binding:"required"`
	User string
}

type UsersPayloadUpdate struct {
	Data UsersDataUpdate `json:"data" binding:"required"`
	User string
}

type UsersDataUpdate struct {
	Param User `json:"param" binding:"required"`
	Body  User `json:"body" binding:"required"`
}

type UserPayloadDelete struct {
	Param []User `json:"param" binding:"required"`
}
