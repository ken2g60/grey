package middlewares

type JwtAuthPayload struct {
	TID  string `json:"tid" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type JwtSessionPayload struct {
	Authorized bool   `json:"authorized" binding:"required"`
	Exp        int    `json:"exp" binding:"required"`
	UserID     string `json:"user_id" binding:"required"`
}
