package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/superz97/go-aggregator/internal/database"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("205")).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255"))

	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Italic(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120"))
)

type tuiModel struct {
	posts    []database.GetPostsForTUIRow
	cursor   int
	offset   int
	viewport int
	detail   bool
	err      error
	status   string
}

func newTUIModel(posts []database.GetPostsForTUIRow) tuiModel {
	return tuiModel{
		posts:    posts,
		viewport: 10,
	}
}

func (m tuiModel) Init() tea.Cmd {
	return nil
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.detail {
		switch keyMsg.String() {
		case "q", "esc":
			m.detail = false
			m.status = ""
		case "o":
			post := m.posts[m.cursor]
			if err := openInBrowser(post.Url); err != nil {
				m.status = fmt.Sprintf("could not open browser: %v", err)
			} else {
				m.status = "opened in browser"
			}
		case "ctrl+c":
			return m, tea.Quit
		}
		return m, nil
	}

	switch keyMsg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			if m.cursor < m.offset {
				m.offset = m.cursor
			}
		}
	case "down", "j":
		if m.cursor < len(m.posts)-1 {
			m.cursor++
			if m.cursor >= m.offset+m.viewport {
				m.offset = m.cursor - m.viewport + 1
			}
		}
	case "enter":
		if len(m.posts) > 0 {
			m.detail = true
			m.status = ""
		}
	case "o":
		if len(m.posts) > 0 {
			post := m.posts[m.cursor]
			if err := openInBrowser(post.Url); err != nil {
				m.status = fmt.Sprintf("could not open browser: %v", err)
			} else {
				m.status = "opened in browser"
			}
		}
	}
	return m, nil
}

func (m tuiModel) View() string {
	if len(m.posts) == 0 {
		return "No posts found. Follow a feed and run `gator agg` first.\n\nq: quit\n"
	}

	if m.detail {
		return m.detailView()
	}
	return m.listView()
}

func (m tuiModel) listView() string {
	var b strings.Builder

	b.WriteString(headerStyle.Render("Posts") + helpStyle.Render("  (enter: view, o: open in browser, q: quit)") + "\n\n")

	end := m.offset + m.viewport
	if end > len(m.posts) {
		end = len(m.posts)
	}

	for i := m.offset; i < end; i++ {
		post := m.posts[i]
		line := "  " + post.Title
		if i == m.cursor {
			line = selectedStyle.Render("> " + post.Title)
		}
		fmt.Fprintln(&b, line)
	}

	if m.status != "" {
		fmt.Fprintf(&b, "\n%s\n", statusStyle.Render(m.status))
	}

	return b.String()
}

func (m tuiModel) detailView() string {
	post := m.posts[m.cursor]

	published := "unknown date"
	if post.PublishedAt.Valid {
		published = post.PublishedAt.Time.Format("2006-01-02")
	}

	description := "(no description)"
	if post.Description.Valid && post.Description.String != "" {
		description = stripHTML(post.Description.String)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s\n%s\n\n%s\n\n%s\n\n",
		titleStyle.Render(post.Title),
		metaStyle.Render(post.FeedName+" · "+published),
		description,
		post.Url,
	)
	b.WriteString(helpStyle.Render("esc/q: back, o: open in browser") + "\n")

	if m.status != "" {
		fmt.Fprintf(&b, "\n%s\n", m.status)
	}

	return b.String()
}
