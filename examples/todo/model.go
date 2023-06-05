package todo

type Todo struct {
	ID        uint   `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}
