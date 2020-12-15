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

type entryUpdateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type entryAddTagsRequest struct {
	Tags []string `json:"tags"`
}

type entryRemoveTagRequest struct {
	Tag string `json:"tag"`
}

type Entry struct {
	Id          string   `json:"id,omitempty"`
	Url         string   `json:"entry_url,omitempty"`
	JournalURL  string   `json:"journal_url,omitempty"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	ContextUrl  string   `json:"context_url,omitempty"`
	ContextType string   `json:"context_type,omitempty"`
	Score       float64  `json:"score,omitempty"`
}

type EntryResultsPage struct {
	TotalResults int     `json:"total_results"`
	Offset       int     `json:"offset,omitempty"`
	NextOffset   int     `json:"next_offset,omitempty"`
	MaxScore     float64 `json:"max_score"`
	Results      []Entry `json:"results"`
}
