package inputs

type CreateWishlistInput struct {
	Title string `json:"title" bind:"min=3,max=100"`
}
