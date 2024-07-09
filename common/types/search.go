package types

type SearchId struct {
	Id uint `json:"id" form:"id" uri:"id" binding:"required,number"`
}
