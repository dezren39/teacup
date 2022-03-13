// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/dirfs"
)

// SetSize sets the size of the filetree.
func (b *Bubble) SetSize(width, height int) {
	horizontal, vertical := bubbleStyle.GetFrameSize()

	b.list.Styles.StatusBar.Width(width - horizontal)
	b.list.SetSize(width-horizontal, height-vertical-lipgloss.Height(b.input.View())-inputStyle.GetVerticalPadding())
}

// SetBorderColor sets the color of the border.
func (b *Bubble) SetBorderColor(color lipgloss.AdaptiveColor) {
	bubbleStyle = bubbleStyle.Copy().BorderForeground(color)
}

// GetSelectedItem returns the currently selected item in the tree.
func (b Bubble) GetSelectedItem() Item {
	selectedDir, ok := b.list.SelectedItem().(Item)
	if ok {
		return selectedDir
	}

	return Item{}
}

// Update handles updating the filetree.
func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case getDirectoryListingMsg:
		if msg != nil {
			cmd = b.list.SetItems(msg)
			cmds = append(cmds, cmd)
		}
	case copyToClipboardMsg:
		return b, b.list.NewStatusMessage(statusMessageInfoStyle(string(msg)))
	case errorMsg:
		return b, b.list.NewStatusMessage(statusMessageErrorStyle(msg.Error()))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, openDirectoryKey):
			if !b.input.Focused() {
				selectedDir := b.GetSelectedItem()
				cmds = append(cmds, getDirectoryListingCmd(selectedDir.FileName, b.showHidden))
			}
		case key.Matches(msg, copyItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully copied file"),
				)

				cmds = append(cmds, tea.Sequentially(
					copyItemCmd(selectedItem.FileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)
			}
		case key.Matches(msg, zipItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully zipped item"),
				)

				cmds = append(cmds, tea.Sequentially(
					zipItemCmd(selectedItem.FileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)
			}
		case key.Matches(msg, unzipItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully unzipped item"),
				)

				cmds = append(cmds, tea.Sequentially(
					unzipItemCmd(selectedItem.FileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)
			}
		case key.Matches(msg, createFileKey):
			if !b.input.Focused() {
				b.input.Focus()
				b.input.Placeholder = "Enter name of new file"
				b.state = createFileState

				return b, textinput.Blink
			}
		case key.Matches(msg, createDirectoryKey):
			if !b.input.Focused() {
				b.input.Focus()
				b.input.Placeholder = "Enter name of new directory"
				b.state = createDirectoryState

				return b, textinput.Blink
			}
		case key.Matches(msg, deleteItemKey):
			if !b.input.Focused() {
				b.input.Focus()
				b.input.Placeholder = "Are you sure you want to delete (y/n)?"
				b.state = deleteItemState

				return b, textinput.Blink
			}
		case key.Matches(msg, toggleHiddenKey):
			if !b.input.Focused() {
				b.showHidden = !b.showHidden
				cmds = append(cmds, getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden))
			}
		case key.Matches(msg, homeShortcutKey):
			if !b.input.Focused() {
				cmds = append(cmds, getDirectoryListingCmd(dirfs.HomeDirectory, b.showHidden))
			}
		case key.Matches(msg, copyToClipboardKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				cmds = append(cmds, copyToClipboardCmd(selectedItem.FileName))
			}
		case key.Matches(msg, escapeKey):
			if b.input.Focused() {
				b.input.Reset()
				b.input.Blur()
				b.state = idleState
			}
		case key.Matches(msg, submitInputKey):
			switch b.state {
			case idleState:
				return b, nil
			case createFileState:
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully created file"),
				)

				cmds = append(cmds, tea.Sequentially(
					createFileCmd(b.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)

				b.input.Blur()
				b.input.Reset()
			case createDirectoryState:
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully created directory"),
				)

				cmds = append(cmds, statusCmd)
				cmds = append(cmds, tea.Sequentially(
					createDirectoryCmd(b.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))

				b.input.Blur()
				b.input.Reset()
			case deleteItemState:
				if strings.ToLower(b.input.Value()) == "y" {
					selectedDir := b.GetSelectedItem()

					statusCmd := b.list.NewStatusMessage(
						statusMessageInfoStyle("Successfully deleted item"),
					)

					cmds = append(cmds, statusCmd)
					cmds = append(cmds, tea.Sequentially(
						deleteItemCmd(selectedDir.FileName),
						getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
					))
				}

				b.input.Blur()
				b.input.Reset()
			}
		}
	}

	switch b.state {
	case idleState:
		b.list, cmd = b.list.Update(msg)
		cmds = append(cmds, cmd)
	case createFileState, createDirectoryState, deleteItemState:
		b.input, cmd = b.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return b, tea.Batch(cmds...)
}
