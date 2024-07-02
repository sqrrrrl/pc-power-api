package api

type LoginCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterCredentials struct {
	Username string `json:"username" binding:"required,printascii,excludesall= ,min=1,max=32"`
	Password string `json:"password" binding:"required,min=8,max=128"`
	Confirm  string `json:"confirm" binding:"required,eqfield=Password"`
}
