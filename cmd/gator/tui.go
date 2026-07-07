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

	frameStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(0, 1)

	footerAccentStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("205")).
				Foreground(lipgloss.Color("0")).
				Bold(true).
				Padding(0, 1)

	footerSegmentStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("236")).
				Foreground(lipgloss.Color("250")).
				Padding(0, 1)
)

type tuiModel struct {
	posts    []database.GetPostsForTUIRow
	cursor   int
	offset   int
	viewport int
	detail   bool
	err      error
	status   string
	width    int
	height   int
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport = m.height - 8
		if m.viewport < 3 {
			m.viewport = 3
		}
		if m.cursor >= m.offset+m.viewport {
			m.offset = m.cursor - m.viewport + 1
		}
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m tuiModel) handleKey(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

func (m tuiModel) contentWidth() int {
	width := m.width
	if width <= 0 {
		width = 60
	}
	cw := width - 4
	if cw < 20 {
		cw = 20
	}
	return cw
}

func (m tuiModel) footer(accent, help string) string {
	right := help
	if m.status != "" {
		right = help + " " + m.status
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, footerAccentStyle.Render(accent), footerSegmentStyle.Render(right))
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
		b.WriteString(line + "\n")
	}

	header := headerStyle.Render("Posts")
	box := frameStyle.Width(m.contentWidth()).Render(strings.TrimRight(b.String(), "\n"))
	footer := m.footer(
		fmt.Sprintf(" POST %d/%d ", m.cursor+1, len(m.posts)),
		"enter view · o open · q quit",
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, box, footer)
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

	content := fmt.Sprintf("%s\n%s\n\n%s\n\n%s",
		titleStyle.Render(post.Title),
		metaStyle.Render(post.FeedName+" · "+published),
		description,
		post.Url,
	)

	header := headerStyle.Render("Post detail")
	box := frameStyle.Width(m.contentWidth()).Render(content)
	footer := m.footer(
		fmt.Sprintf(" POST %d/%d ", m.cursor+1, len(m.posts)),
		"esc/q back · o open",
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, box, footer)
}
