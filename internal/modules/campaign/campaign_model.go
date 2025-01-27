package campaign

type Campaign struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GoalAmount  int    `json:"goal_amount"`
}
