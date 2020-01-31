package model

import (
	"fmt"
	"time"
	"log"

	"app/shared/database"
)

// *****************************************************************************
// Post
// *****************************************************************************

// Post struct contains the information for each post
type Post struct {
	ID        uint32 `db:"id"`
	Content   string `db:"content"`
	UserID    uint32 `db:"user_id"`
	Files     []Upload
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Deleted   uint8     `db:"deleted"`
}

// PostByID gets post by ID
func PostByID(postID string, userID string) (Post, error) {
	result := Post{}
	err := database.SQL.Get(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM post WHERE id = ? AND user_id = ? LIMIT 1", postID, userID)
	return result, StandardizeError(err)
}

// PostsByUserID gets all posts for a user
func PostsByUserID(userID string) ([]Post, error) {
	var result []Post
	err := database.SQL.Select(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM post WHERE user_id = ?", userID)
	// Get all uploads for a post
	for _, r := range result {
		uploads, err := UploadsByPostID(r.ID)
		if err != nil {
			log.Println(err)
		}
		r.UploadsGET(uploads)
	}
	log.Println(result)
	return result, StandardizeError(err)
}

// PostCreate creates a post and returns its id
func PostCreate(content string, userID string) (string, error, error) {
	_, e := database.SQL.Exec("INSERT INTO post (content, user_id) VALUES (?,?)", content, userID)
	result := Post{}
	err := database.SQL.Get(&result, "SELECT id FROM post WHERE content = ? AND user_id = ? LIMIT 1", content, userID)
	// convert uint32 to string
	id := fmt.Sprint(result.ID)
	return id, StandardizeError(e), StandardizeError(err)
}

// PostUpdate updates a post
func PostUpdate(content string, userID string, postID string) error {
	_, err := database.SQL.Exec("UPDATE post SET content=? WHERE user_id = ? AND id = ? LIMIT 1", content, userID, postID)
	return StandardizeError(err)
}

// PostDelete deletes a post
func PostDelete(postID string, userID string) error {
	_, err := database.SQL.Exec("DELETE FROM post WHERE id = ? AND user_id = ?", postID, userID)
	return StandardizeError(err)
}
