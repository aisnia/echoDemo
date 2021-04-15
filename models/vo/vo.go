package vo

type RegisterReq struct {
	Name     string `json:"name"  validate:"required"`
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password"  validate:"required"`
	Token    string `json:"token"  validate:"required"`
	Code     string `json:"code" validate:"required"`
}
