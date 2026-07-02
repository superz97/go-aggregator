package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	"github.com/superz97/go-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	fs := flag.NewFlagSet("browse", flag.ContinueOnError)

	limit := fs.Int("limit", 2, "limit")
	page := fs.Int("page", 1, "page")
	sort := fs.String("sort", "published_at", "sort")
	order := fs.String("order", "desc", "order")
	feed := fs.String("feed", "", "")

	if err := fs.Parse(cmd.args); err != nil {
		return err
	}

	if *page < 1 {
		return fmt.Errorf("invalid page %d", *page)
	}
	offset := (*page - 1) * *limit

	switch *sort {
	case "published_at", "title":
	default:
		return fmt.Errorf("invalid sort %q", *sort)
	}

	switch *order {
	case "asc", "desc":
	default:
		return fmt.Errorf("invalid order %q", *order)
	}

	feedName := sql.NullString{
		String: *feed,
		Valid:  *feed != "",
	}

	var (
		posts []database.Post
		err   error
	)

	key := fmt.Sprintf("%s:%s", *sort, *order)
	limit32 := int32(*limit)
	offset32 := int32(offset)

	switch key {
	case "published_at:asc":
		posts, err = s.db.GetPostsForUserByPublishedAtAsc(
			context.Background(),
			database.GetPostsForUserByPublishedAtAscParams{
				UserID:   user.ID,
				Limit:    limit32,
				Offset:   offset32,
				FeedName: feedName,
			},
		)
	case "published_at:desc":
		posts, err = s.db.GetPostsForUserByPublishedAtDesc(
			context.Background(),
			database.GetPostsForUserByPublishedAtDescParams{
				UserID:   user.ID,
				Limit:    limit32,
				Offset:   offset32,
				FeedName: feedName,
			},
		)
	case "title:asc":
		posts, err = s.db.GetPostsForUserByTitleAsc(
			context.Background(),
			database.GetPostsForUserByTitleAscParams{
				UserID:   user.ID,
				Limit:    limit32,
				Offset:   offset32,
				FeedName: feedName,
			},
		)
	case "title:desc":
		posts, err = s.db.GetPostsForUserByTitleDesc(
			context.Background(),
			database.GetPostsForUserByTitleDescParams{
				UserID:   user.ID,
				Limit:    limit32,
				Offset:   offset32,
				FeedName: feedName,
			},
		)
	}

	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("%s\n%s\n\n", post.Title, post.Url)
	}

	return nil
}
