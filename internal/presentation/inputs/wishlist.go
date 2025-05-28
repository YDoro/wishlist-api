package inputs

type CreateWishlistInput struct {
	Title string `json:"title" bind:"min=3,max=100"`
}

type UpdateWishlistInput struct {
	Title string   `json:"title"`
	Items []string `json:"items"`
}
