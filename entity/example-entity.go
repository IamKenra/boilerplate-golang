package entity

import "time"

type AssetMaintenanceKey struct {
	IdAsset  uint64 `json:"-"`
	IdStatus uint64 `json:"-"`
}

type AssetMaintenance struct {
	AssetMaintenanceKey

	UUIDAsset  string `json:"uuid_asset,omitempty" form:"uuid_asset"`
	UUIDStatus string `json:"uuid_status,omitempty" form:"uuid_status"`

	Judul     string     `json:"judul,omitempty" from:"judul"`
	Catatan   string     `json:"catatan,omitempty" from:"catatan"`
	Tgl       *time.Time `json:"tgl,omitempty" from:"tgl"`
	Eksekutor string     `json:"eksekutor,omitempty" from:"eksekutor"`
	FlagAktif int        `json:"flag_aktif,omitempty" from:"flag_aktif"`
}
