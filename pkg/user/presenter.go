package user

type UserPresenter struct {
	ID   string `json:"id"`
	Name string `json:"name" binding:"required"`
}
