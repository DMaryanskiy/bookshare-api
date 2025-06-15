package books

type Handler struct {}

func NewHandler() *Handler {
	return &Handler{}
}

type BookInput struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author"`
	Description string `json:"description"`
}
