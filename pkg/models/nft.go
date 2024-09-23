package models

type NFT struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
	Price string `json:"price"`
	ImageURL string `json:"image_url"`
}
// //можно добавить ещё много разного

// Description string `json:"description"`

// CurrentBid string `json:"current_bid"`

//Информация об истории продаж или предыдущих владельцах NFTили Может быть списком транзакций
// alesHistory []Sale `json:"sales_history"`
// type Sale struct {
//     Seller    string `json:"seller"`
//     Buyer     string `json:"buyer"`
//     Price     string `json:"price"`
//     Timestamp string `json:"timestamp"`
// }
