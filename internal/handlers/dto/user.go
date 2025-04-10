package dto

const (
	RoleEmployee  = "employee"
	RoleModerator = "moderator"
)

type UserRegisterRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserRegisterResponseDto struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type DummyLoginRequestDto struct {
	Role string `json:"role"`
}

type UserLoginRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponseDto struct {
	Token string `json:"token"`
}
