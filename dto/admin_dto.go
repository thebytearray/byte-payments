package dto

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Admin AdminInfo `json:"admin"`
}

type AdminInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CreatePlanRequest struct {
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	PriceUSD     float64 `json:"price_usd" validate:"required,gt=0"`
	DurationDays int64   `json:"duration_days" validate:"required,gt=0"`
}

type UpdatePlanRequest struct {
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	PriceUSD     float64 `json:"price_usd" validate:"required,gt=0"`
	DurationDays int64   `json:"duration_days" validate:"required,gt=0"`
}

type CreateCurrencyRequest struct {
	Code         string `json:"code" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Network      string `json:"network" validate:"required"`
	IsToken      bool   `json:"is_token"`
	ContractAddr string `json:"contract_addr"`
	Enabled      bool   `json:"enabled"`
}

type UpdateCurrencyRequest struct {
	Code         string `json:"code" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Network      string `json:"network" validate:"required"`
	IsToken      bool   `json:"is_token"`
	ContractAddr string `json:"contract_addr"`
	Enabled      bool   `json:"enabled"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}