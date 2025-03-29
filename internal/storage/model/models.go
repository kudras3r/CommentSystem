package model

type Comment struct {
	ID        string     `json:"id" db:"id"`
	AuthorID  string     `json:"authorID" db:"author_id"`
	Content   string     `json:"content" db:"content"`
	CreatedAt string     `json:"createdAt" db:"created_at"`
	Rating    int32      `json:"rating" db:"rating"`
	PostID    string     `json:"postID" db:"post_id"`
	ParentID  *string    `json:"parentID,omitempty" db:"parent_id"`
	Children  []*Comment `json:"children"`
}

type Post struct {
	ID         string     `json:"id" db:"id"`
	AuthorID   string     `json:"authorID" db:"author_id"`
	Title      string     `json:"title" db:"title"`
	Content    string     `json:"content" db:"content"`
	AllowComms bool       `json:"allowComms" db:"allow_comms"`
	CreatedAt  string     `json:"createdAt" db:"created_at"`
	Rating     int32      `json:"rating" db:"rating"`
	Comments   []*Comment `json:"comments,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type Subscription struct {
}
