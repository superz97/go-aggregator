package main

import (
	"aggregator/internal/database"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login requires a username argument")
	}

	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user %s does not exist", name)
	}

	s.cfg.SetUser(name)
	fmt.Printf("User has been set to: %s\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("register requires at least one argument")
	}

	name := cmd.args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("error: user already exist or could not be created: %w", err)
	}

	s.cfg.SetUser(name)
	fmt.Printf("user created: %+v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}
	fmt.Println("Database reset successfully")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
