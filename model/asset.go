package model

func init() {
	Register(&Asset{})
}

type Asset struct {
	Base
	ID            string `json:"id"`
	Class         string `json:"class"`
	Exchange      string `json:"exchange"`
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	Tradable      bool   `json:"tradable"`
	Marginable    bool   `json:"marginable"`
	Shortable     bool   `json:"shortable"`
	EasyToBorrow  bool   `json:"easy_to_borrow"`
	Fractionable  bool   `json:"fractionable"`
	IsWatchlisted bool   `json:"is_watchlisted"`
}

type AssetsRepo interface {
	CreateOrUpdate(*Asset) (*Asset, error)
	UpdateAsset(*Asset) error
	Search(string) ([]Asset, error)
}
