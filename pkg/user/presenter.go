package user

type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

type UserPresenter struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
