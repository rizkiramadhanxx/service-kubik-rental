package entity

const (
	ModuleUser        = "user"
	ModuleSetting     = "setting"
	ModuleRole        = "role"
	ModuleTransaction = "transaction"
	ModuleCategory    = "category"
	ModuleProduct     = "product"
	ModuleDevice      = "device"
	ModuleBilling     = "billing"
	ModuleMember      = "member"
	ModulePackage     = "package"
	ModulePOS         = "pos"
)

var AllModules = []string{
	ModuleUser,
	ModuleSetting,
	ModuleRole,
	ModuleTransaction,
	ModuleCategory,
	ModuleProduct,
	ModuleDevice,
	ModuleBilling,
	ModuleMember,
	ModulePackage,
	ModulePOS,
}

type Role struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Name    string `json:"name" validate:"required" gorm:"unique"`
	Modules string `json:"module" validate:"required"`
}
