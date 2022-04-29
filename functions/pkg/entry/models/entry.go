package models

const (
	EntryModelName = "entry"
)

type Entry struct {
	Hash  string `json:"hash" binding:"required,alphanum,len=64"`
	Name  string `json:"name" binding:"required,alphanum,min=3,max=32"`
	Score uint   `json:"score" binding:"required,gte=0,lte=9999999"`
}
