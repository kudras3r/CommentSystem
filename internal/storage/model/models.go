package model

type CommEdge struct {
	Node   *Comment `json:"node"`
	Cursor string   `json:"cursor"`
}

type Comment struct {
	ID        string  `json:"id" db:"id"`
	AuthorID  string  `json:"authorID" db:"author_id"`
	Content   string  `json:"content" db:"content"`
	CreatedAt string  `json:"createdAt" db:"created_at"`
	Rating    int32   `json:"rating" db:"rating"`
	PostID    string  `json:"postID" db:"post_id"`
	ParentID  *string `json:"parentID,omitempty" db:"parent_id"`
}

type CommentsConnection struct {
	Edges    []*CommEdge `json:"edges"`
	PageInfo *PageInfo   `json:"pageInfo"`
}

type Mutation struct {
}

type PageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

type Post struct {
	ID         string              `json:"id" db:"id"`
	AuthorID   string              `json:"authorID" db:"author_id"`
	Title      string              `json:"title" db:"title"`
	Content    string              `json:"content" db:"content"`
	AllowComms bool                `json:"allowComms" db:"allow_comms"`
	CreatedAt  string              `json:"createdAt" db:"created_at"`
	Rating     int32               `json:"rating" db:"rating"`
	Comments   *CommentsConnection `json:"comments" db:"comments"`
}

type Query struct {
}

type Subscription struct {
}

type User struct {
	ID       string     `json:"id" db:"id"`
	Username string     `json:"username" db:"username"`
	Email    *string    `json:"email,omitempty" db:"email"`
	Posts    []*Post    `json:"posts"`
	Comments []*Comment `json:"comments"`
}
