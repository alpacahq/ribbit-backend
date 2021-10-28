package model

type PlaidAuthToken struct {
	LinkToken string `json:"link_token"`
}

type AccessToken struct {
	ID          int    `json:"id"`
	PublicToken string `json:"public_token"`
	ItemID      string `json:"item_id"`
}
