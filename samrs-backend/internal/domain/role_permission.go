package domain

type RolePermission struct {
	RoleID       int `gorm:"primaryKey" json:"role_id"`
	PermissionID int `gorm:"primaryKey" json:"permission_id"`
}
