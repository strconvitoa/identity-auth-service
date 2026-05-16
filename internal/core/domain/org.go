package domain

type Org struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	TradeName string `json:"tradename"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	URL       string `json:"url"`
	Image     string `json:"image"`
}
