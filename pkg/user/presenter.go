package user

type CreateLocalUserRequest struct {
	Name string `json:"name" binding:"required"`
}

type GetUserRequest struct {
	Uid string `form:"uid" binding:"required"`
}

type UserPresenter struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GoogleUserPresenter struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
