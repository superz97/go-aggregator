package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/superz97/go-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) == 1 {
		parsed, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit %q: %w", cmd.args[0], err)
		}
		limit = parsed
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("%s\n%s\n\n", post.Title, post.Url)
	}
	return nil
}
