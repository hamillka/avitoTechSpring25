package dto

const (
	RoleEmployee  = "employee"
	RoleModerator = "moderator"
)

// UserRegisterRequestDto model info
// @Description Информация о пользователе при регистрации
type UserRegisterRequestDto struct {
	Email    string `json:"email"`    // Почта
	Password string `json:"password"` // Пароль
	Role     string `json:"role"`     // Роль пользователя (employee || moderator)
}

// UserRegisterResponseDto model info
// @Description Информация о пользователе при регистрации
type UserRegisterResponseDto struct {
	Id    string `json:"id"`    // Идентификатор
	Email string `json:"email"` // Почта
	Role  string `json:"role"`  // Роль пользователя (employee || moderator)
}

// DummyLoginRequestDto model info
// @Description Информация о пользователе при упрощенном входе
type DummyLoginRequestDto struct {
	Role string `json:"role"` // Желаемая роль (employee || moderator)
}

// UserLoginRequestDto model info
// @Description Информация о пользователе при входе в систему
type UserLoginRequestDto struct {
	Email    string `json:"email"`    // Почта
	Password string `json:"password"` // Пароль
}

// UserLoginResponseDto model info
// @Description Информация о пользователе при входе в систему
type UserLoginResponseDto struct {
	Token string `json:"token"` // JWT токен
}
