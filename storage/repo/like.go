package repo

type Like struct {
	ID int64
	PostID int64
	UserID int64
	Status bool
}

type LikesDislikesCountsResult struct {
	LikesCount int64
	DislikesCount int64
}

type LikeStorageI interface {
	Create(l *Like) (*Like, error)
	Get(userID, postID int64) (*Like, error)
	GetLikesDislikeCount(postID int64) (*LikesDislikesCountsResult, error)
}