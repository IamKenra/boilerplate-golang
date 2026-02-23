package entity

type Identity struct {
	Name         string `json:"name,omitempty"`
	Alias        string `json:"alias,omitempty"`
	Address      string `json:"address,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	Phone        string `json:"phone,omitempty"`
}
