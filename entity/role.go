package entity

const (
	ModuleUser        = "user"
	ModuleSetting     = "setting"
	ModuleRole        = "role"
	ModuleTransaction = "transaction"
)

var AllModules = []string{
	ModuleUser,
	ModuleSetting,
	ModuleRole,
	ModuleTransaction,
}

type Role struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Name    string `json:"name" validate:"required" gorm:"unique"`
	Modules string `json:"module" validate:"required"`
}
