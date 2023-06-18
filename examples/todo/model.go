package todo

// Represents a Todo item.
type Todo struct {
	// Id of the todo item
	ID        uint   `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"` // whether the todo item is completed or not
}
