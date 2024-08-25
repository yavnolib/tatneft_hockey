package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"tat_hockey_pack/internal/models"
	"time"
)

type Post struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *Post {
	return &Post{
		db: db,
	}
}

func (p *Post) GetAllPosts() ([]models.PostPreview, error) {
	rows, err := p.db.Query(context.Background(), `SELECT id, preview, creator_id, created_at FROM posts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostPreview
	for rows.Next() {
		var post models.PostPreview
		err := rows.Scan(&post.ID, &post.Preview, &post.CreatorID, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *Post) GetPostByID(id int64) (models.Post, error) {
	var post models.Post

	row := p.db.QueryRow(context.Background(), `
		SELECT id, title, preview, video_id, creator_id, created_at
		FROM posts
		WHERE id = $1`, id)

	err := row.Scan(&post.ID, &post.Title, &post.Preview, &post.VideoID, &post.CreatorID, &post.CreatedAt)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (p *Post) CreatePost(ctx context.Context, videoID int64, title, preview string, creatorID int64) (int, error) {
	query := `
        INSERT INTO posts (title, preview, video_id, creator_id, created_at) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id;
    `

	var postID int
	err := p.db.QueryRow(ctx, query, title, preview, videoID, creatorID, time.Now()).Scan(&postID)
	if err != nil {
		return -1, fmt.Errorf("failed to create post: %w", err)
	}

	fmt.Printf("Post created with ID: %d\n", postID)
	return postID, nil
}

func (p *Post) SaveGIF(ctx context.Context, path, className string, postID int) (int, error) {
	query := `
        INSERT INTO gifs (path, class_name, post_id)
        VALUES ($1, $2, $3)
        RETURNING id;
    `

	var gifID int
	err := p.db.QueryRow(ctx, query, path, className, postID).Scan(&gifID)
	if err != nil {
		return -1, fmt.Errorf("failed to save GIF: %w", err)
	}

	fmt.Printf("GIF saved with ID: %d\n", gifID)
	return gifID, nil
}
