package transaction

type CheckoutRequest struct {
	CartID    uint    `json:"cart_id" validate:"required"`
	MemberID  *uint   `json:"member_id,omitempty"`
	BuyerName *string `json:"buyer_name,omitempty"`
}

type TransactionQuery struct {
	StartDate string `query:"start_date"` // format: "2006-01-02"
	EndDate   string `query:"end_date"`
	Type      string `query:"type"`    // optional: product, billing, mixed
	Keyword   string `query:"keyword"` // search by buyer_name
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
}
