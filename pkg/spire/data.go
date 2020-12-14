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

type EntryContext struct {
	ContextType string
	ContextID   string
	ContextURL  string
}

type entryCreateRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	ContextType string   `json:"context_type,omitempty"`
	ContextURL  string   `json:"context_url,omitempty"`
	ContextID   string   `json:"context_id,omitempty"`
}

type Entry struct {
	Id          string   `json:"id"`
	JournalURL  string   `json:"journal_url"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	ContextUrl  string   `json:"context_url,omitempty"`
	ContextType string   `json:"context_type,omitempty"`
}

type EntriesList struct {
	Entries []Entry `json:"entries"`
}
