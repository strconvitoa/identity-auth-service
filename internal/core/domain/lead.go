package domain

type Lead struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Issue       string `json:"issue"`
	OrgID       string `json:"org_id"`
	UserID      string `json:"user_id,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Lang        string `json:"lang,omitempty"`
	Status      string `json:"status"`
	Area        string `json:"area"`
	Description string `json:"description"`
	Urgency     string `json:"urgency"`
	Quality     string `json:"quality"`
	Conflict    string `json:"conflict"`
	Retention   string `json:"retention"`
	Created     string `json:"created"`
}
type LeadsResp struct {
	Leads     []Lead `json:"leads"`
	Discarded []Lead `json:"discarded_leads"`
	Pending   []Lead `json:"pending"`
	NewLeads  int    `json:"new_leads"`
}
