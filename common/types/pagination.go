package types

type Pagination struct {
	PageNumber int    `json:"PageNumber" form:"PageNumber" uri:"PageNumber" binding:"validateDefault=1"`
	PageSize   int    `json:"PageSize" form:"PageSize" uri:"PageSize"  binding:"validateDefault=10"`
	Sort       string `json:"sort" form:"Sort" uri:"Sort" binding:"validateDefault=DESC,oneof=ASC DESC" `
}
