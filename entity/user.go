package entity

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required" gorm:"unique"`
	Password string `json:"password,omitempty" validate:"required"` // `omitempty` agar tidak muncul di JSON response
	RoleID   uint   `json:"role_id" validate:"required"`            // untuk input, foreign key ke Role
	Role     Role   `json:"role" gorm:"foreignKey:RoleID"`          // preload relasi
}
