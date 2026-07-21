package dto

type UserLoginResponse struct {
	SessionId string `json:"session_id"`
}

type UserLoginReq struct {
	Username string `json:"username" example:"testusername" validate:"required"`
	Password string `json:"password" example:"testpassword" validate:"required"`
}

type ChangeUsernameReq struct {
	Username string `json:"username" example:"testusername" validate:"required"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password" example:"oldpassword" validate:"required"`
	NewPassword string `json:"new_password" example:"newpassword" validate:"required"`
}

type UserSession struct {
	SessionId string `json:"session_id"`
	UserId    uint   `json:"user_id"`
	Username  string `json:"name"`
	Role      string `json:"role"`
}
