package model

type ActivityItems struct {
	ItemID           int64   `json:"ItemID"`
	ItemName         string  `json:"ItemName,omitempty"`
	ItemDesc         string  `json:"ItemDesc,omitempty"`
	ItemPrice        float64 `json:"ItemPrice,string"`
	ItemSpecialPrice float64 `json:"ItemSpecialPrice,string"`
	PromotionID      int64   `json:"PromotionID,string"`
}

type Items struct {
	ItemID         int64   `json:"item_id"`
	ItemName       string  `json:"item_name"`
	ItemDesc       string  `json:"item_desc"`
	ItemPrice      float64 `json:"item_price"`
	ItemStock      int64   `json:"item_stock"`
	Status         int64   `json:"status"`
	CreatorID      int64   `json:"creator_id"`
	LastModifierID int64   `json:"last_modifier_id"`
}

type Promotions struct {
	PromotionID      int64   `json:"promotion_id"`
	PromotionDesc    string  `json:"promotion_desc"`
	ItemSpecialPrice float64 `json:"item_special_price"`
	StartTime        string  `json:"start_time"`
	EndTime          string  `json:"end_time"`
	CreatorID        int64   `json:"creator_id"`
	LastModifierID   int64   `json:"last_modifier_id"`
	ItemID           int64   `json:"item_id"`
}
