package main

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/superz97/go-aggregator/internal/database"
)

func handlerTUI(s *state, cmd command, user database.User) error {
	posts, err := s.db.GetPostsForTUI(context.Background(), database.GetPostsForTUIParams{
		UserID: user.ID,
		Limit:  100,
	})
	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}

	p := tea.NewProgram(newTUIModel(posts))
	_, err = p.Run()
	if err != nil {
		return fmt.Errorf("could not run tui: %w", err)
	}
	return nil
}
