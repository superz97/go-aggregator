package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/superz97/go-aggregator/internal/database"
)

func handlerLike(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("like requires a post url argument")
	}

	post, err := s.db.GetPostByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not find post with that url: %w", err)
	}

	like, err := s.db.CreatePostLike(context.Background(), database.CreatePostLikeParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		PostID:    post.ID,
	})
	if err != nil {
		return fmt.Errorf("could not like post: %w", err)
	}

	fmt.Printf("Liked post: %s\n", like.PostTitle)
	return nil
}

func handlerUnlike(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("unlike requires a post url argument")
	}

	post, err := s.db.GetPostByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not find post with that url: %w", err)
	}

	err = s.db.DeletePostLike(context.Background(), database.DeletePostLikeParams{
		UserID: user.ID,
		PostID: post.ID,
	})
	if err != nil {
		return fmt.Errorf("could not unlike post: %w", err)
	}

	return nil
}

func handlerLikes(s *state, cmd command, user database.User) error {
	posts, err := s.db.GetLikedPostsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get liked posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("%s\n%s\n\n", post.Title, post.Url)
	}
	return nil
}
