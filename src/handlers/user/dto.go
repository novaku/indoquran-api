package user

type formPost struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	FacebookID string `json:"facebook_id"`
}

type formLogout struct {
	UserID string `json:"user_id" binding:"required"`
}
