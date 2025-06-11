package dto

type UserRegisterOutput struct {
	Body UserRegisterResponse
}

type UserReturn struct {
	ID        int64  `json:"id" example:"687627953"`          // Telegram ID
	Username  string `json:"username" example:"warl0rdd"`     // Telegram username
	Firstname string `json:"first_name" example:"Linuxfight"` // Telegram first name
	Lastname  string `json:"last_name" example:"Olukhovich"`  // Telegram last name
	PhotoUrl  string `json:"photo_url"`                       // Telegram photo URL
}

type UserRegisterResponse struct {
	User   UserReturn `json:"user"`   // User object
	Tokens AuthTokens `json:"tokens"` // Two JWT tokens: Access token and Refresh token
}

type UserLoginInput struct {
	Body UserLogin
}

type UserLogin struct {
	InitDataRaw string `json:"init_data_raw" validate:"required,init_data_raw" example:"user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%20%2B%20-%20%3F%20%5C%2F%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F4FPEE4tmP3ATHa57u6MqTDih13LTOiMoKoLDRG4PnSA.svg%22%7D&chat_instance=8134722200314281151&chat_type=private&auth_date=1733509682&signature=TYJxVcisqbWjtodPepiJ6ghziUL94-KNpG8Pau-X7oNNLNBM72APCpi_RKiUlBvcqo5L-LAxIc3dnTzcZX_PDg&hash=a433d8f9847bd6addcc563bff7cc82c89e97ea0d90c11fe5729cae6796a36d73"` // init data from telegram // User's password
}
