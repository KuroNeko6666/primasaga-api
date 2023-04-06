package form

type Post struct {
	Caption string `json:"caption" form:"caption" validate:"required"`
}
