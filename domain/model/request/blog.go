package request_model

type BlogRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var AllowedMimeTypes = []string{
	"image/jpeg", "image/png", "image/jpg", // images
}
