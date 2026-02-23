package entity

type UserKey struct {
	IDIdentity uint64 `json:"-" form:"id_identity"`
}
type User struct {
	UserKey
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Username     string `json:"username,omitempty" form:"username"`
	Password     string `json:"password,omitempty" form:"password"`
	Authority    string `json:"authority,omitempty" form:"authority"`
}
