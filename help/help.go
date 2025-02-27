// Package help implements a help bubble which can be used
// to display help information such as keymaps.
package help

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 1
	keyWidth = 12
)

type TitleColor struct {
	Background lipgloss.AdaptiveColor
	Foreground lipgloss.AdaptiveColor
}

// Entry represents a single entry in the help bubble.
type Entry struct {
	Key         string
	Description string
}

// Model represents the properties of a help bubble.
type Model struct {
	Viewport    viewport.Model
	Entries     []Entry
	BorderColor lipgloss.AdaptiveColor
	Title       string
	TitleColor  TitleColor
	Active      bool
	Borderless  bool
}

// generateHelpScreen generates the help text based on the title and entries.
func generateHelpScreen(title string, titleColor TitleColor, entries []Entry, width, height int) string {
	helpScreen := ""

	for _, content := range entries {
		keyText := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"}).
			Width(keyWidth).
			Render(content.Key)

		descriptionText := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"}).
			Render(content.Description)

		row := lipgloss.JoinHorizontal(lipgloss.Top, keyText, descriptionText)
		helpScreen += fmt.Sprintf("%s\n", row)
	}

	titleText := lipgloss.NewStyle().Bold(true).
		Background(titleColor.Background).
		Foreground(titleColor.Foreground).
		Border(lipgloss.NormalBorder()).
		Padding(0, 1).
		Italic(true).
		BorderBottom(true).
		BorderTop(false).
		BorderRight(false).
		BorderLeft(false).
		Render(title)

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			titleText,
			helpScreen,
		))
}

// New creates a new instance of a help bubble.
func New(
	active, borderless bool,
	title string,
	titleColor TitleColor,
	borderColor lipgloss.AdaptiveColor,
	entries []Entry,
) Model {
	viewPort := viewport.New(0, 0)
	border := lipgloss.NormalBorder()

	if borderless {
		border = lipgloss.HiddenBorder()
	}

	viewPort.Style = lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		BorderForeground(borderColor)

	viewPort.SetContent(generateHelpScreen(title, titleColor, entries, 0, 0))

	return Model{
		Viewport:    viewPort,
		Entries:     entries,
		Title:       title,
		Active:      active,
		Borderless:  borderless,
		BorderColor: borderColor,
		TitleColor:  titleColor,
	}
}

// SetSize sets the size of the help bubble.
func (m *Model) SetSize(w, h int) {
	m.Viewport.Width = w
	m.Viewport.Height = h

	m.Viewport.SetContent(generateHelpScreen(m.Title, m.TitleColor, m.Entries, m.Viewport.Width, m.Viewport.Height))
}

// SetBorderColor sets the current color of the border.
func (m *Model) SetBorderColor(color lipgloss.AdaptiveColor) {
	m.BorderColor = color
}

// SetIsActive sets if the bubble is currently active.
func (m *Model) SetIsActive(active bool) {
	m.Active = active
}

// GotoTop jumps to the top of the viewport.
func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

// SetTitleColor sets the color of the title.
func (m *Model) SetTitleColor(color TitleColor) {
	m.TitleColor = color

	m.Viewport.SetContent(generateHelpScreen(m.Title, m.TitleColor, m.Entries, m.Viewport.Width, m.Viewport.Height))
}

// SetBorderless sets weather or not to show the border.
func (m *Model) SetBorderless(borderless bool) {
	m.Borderless = borderless
}

// Update handles UI interactions with the help bubble.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if m.Active {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the help bubble.
func (m Model) View() string {
	border := lipgloss.NormalBorder()

	if m.Borderless {
		border = lipgloss.HiddenBorder()
	}

	m.Viewport.Style = lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		BorderForeground(m.BorderColor)

	return m.Viewport.View()
}
