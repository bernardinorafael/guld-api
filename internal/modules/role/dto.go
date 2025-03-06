package role

type ManagePermissionsProps struct {
	Permissions []string `json:"permissions"`
}

type UpdateRoleDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
