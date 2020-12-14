package spire

type journalCreateRequest struct {
	Name string `json:"name"`
}

type Journal struct {
	Id           string   `json:"id"`
	BugoutUserID string   `json:"bugout_user_id"`
	HolderIDs    []string `json:"holder_ids"`
	Name         string   `json:"name"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

type JournalsList struct {
	Journals []Journal `json:"journals"`
}
