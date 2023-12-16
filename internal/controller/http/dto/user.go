package dto

type ChangePasswordDTO struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type TokensDTO struct {
	AccessToken  string `json:"acces_token"`
	RefreshToken string `json:"refresh_token"`
}
