package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/superz97/go-aggregator/internal/database"
)

func handlerBookmark(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("bookmark requires a post url argument")
	}

	post, err := s.db.GetPostByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not find post with that url: %w", err)
	}

	bookmark, err := s.db.CreatePostBookmark(context.Background(), database.CreatePostBookmarkParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		PostID:    post.ID,
	})
	if err != nil {
		return fmt.Errorf("could not bookmark post: %w", err)
	}

	fmt.Printf("Bookmarked post: %s\n", bookmark.PostTitle)
	return nil
}

func handlerUnbookmark(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("unbookmark requires a post url argument")
	}

	post, err := s.db.GetPostByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not find post with that url: %w", err)
	}

	err = s.db.DeletePostBookmark(context.Background(), database.DeletePostBookmarkParams{
		UserID: user.ID,
		PostID: post.ID,
	})
	if err != nil {
		return fmt.Errorf("could not unbookmark post: %w", err)
	}

	return nil
}

func handlerBookmarks(s *state, cmd command, user database.User) error {
	posts, err := s.db.GetBookmarkedPostsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get bookmarked posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("%s\n%s\n\n", post.Title, post.Url)
	}
	return nil
}
