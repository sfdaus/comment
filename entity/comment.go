package entity

type Comment struct {
	ID        string `json:"id"`
	ThreadID  string `json:"thread_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	IsActive  *bool  `json:"is_active"`
	CreatedBy string `json:"created_by"`
	CreatedAt int64  `json:"created_at"`
	UpdatedBy string `json:"updated_by"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}
