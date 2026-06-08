package domain

type Intake struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Issue   string `json:"issue"`
	OrgID   string `json:"org_id"`
	UserID  string `json:"user_id,omitempty"`
	Summary string `json:"summary,omitempty"`
	Lang    string `json:"lang,omitempty"`
}
