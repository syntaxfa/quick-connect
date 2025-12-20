package postgres

import (
	"context"
	"time"

	"github.com/syntaxfa/quick-connect/app/storyapp/service"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const querySaveStory = `INSERT INTO stories (id, media_file_id, title, caption, link_url, link_text,
duration_seconds, is_active, view_count, publish_at, expires_at, creator_id,
created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING id, is_active, view_count, created_at, updated_at;
`

func (d *DB) SaveStory(ctx context.Context, req service.AddStoryRequest) (service.Story, error) {
	const op = "repository.postgres.create.SaveStory"

	var nullable nullableFields

	createdAt := time.Now()
	// Default values for insertion
	viewCount := 0

	story := service.Story{
		MediaFileID:     req.MediaFileID,
		Title:           req.Title,
		Caption:         req.Caption,
		LinkURL:         req.LinkURL,
		LinkText:        req.LinkText,
		DurationSeconds: req.DurationSeconds,
		PublishAt:       req.PublishAt,
		ExpiresAt:       req.ExpiresAt,
		CreatorID:       req.CreatorID,
	}

	// 4. Execute Query
	err := d.conn.Conn().QueryRow(ctx, querySaveStory, req.ID, req.MediaFileID, &nullable.Title,
		&nullable.Caption, &nullable.LinkURL, &nullable.LinkText, req.DurationSeconds, true,
		viewCount, req.PublishAt, req.ExpiresAt, req.CreatorID, createdAt, createdAt,
	).Scan(&story.ID, &story.IsActive, &viewCount, &story.CreatedAt, &story.CreatedAt)

	if err != nil {
		return service.Story{}, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected)
	}

	return story, nil
}
