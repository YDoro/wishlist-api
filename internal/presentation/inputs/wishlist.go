package inputs

type CreateWishlistInput struct {
	Title string `json:"title" bind:"min=3,max=100"`
}

type ProdcutToWishlistInput struct {
	ProductId string `json:"productId" bind:"required"`
}
