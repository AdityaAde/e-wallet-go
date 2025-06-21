package dto

type TopupReq struct {
	Amount float64 `json:"amount"`
	UserID int64   `json:"-"`
}
