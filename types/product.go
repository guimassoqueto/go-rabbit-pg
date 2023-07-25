package types

type Product struct {
	Id string
	Title string
	Category string
	Reviews int
	Free_Shipping bool
	Image_Url string
	Discount int
	Price float32
	Previous_Price float32
}