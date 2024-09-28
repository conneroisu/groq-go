// Package main demonstrates how to use groq-go to create a chat application
// using the groq api accessable through the terminal.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/conneroisu/groq-go"
)

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithMouseAllMotion(),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	groqClient  *groq.Client
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

const (
	width  = 500
	height = 10
)

func initialModel() model {
	groqClient, err := groq.NewClient(
		os.Getenv("GROQ_KEY"),
	)
	if err != nil {
		return model{
			err: errMsg(err),
		}
	}
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.SetWidth(500)
	ta.SetHeight(10)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(width, height)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		groqClient:  groqClient,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyDown:
			m.textarea.SetValue(m.textarea.Value() + "\n")
			m.viewport.LineDown(1)
			m.viewport.TotalLineCount()
		case tea.KeyEnter:
			message := m.textarea.Value()
			if strings.TrimSpace(m.textarea.Value()) == "" {
				break
			}
			m.messages = append(
				m.messages,
				m.senderStyle.Render("You: ")+message,
			)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			re, err := m.groqClient.CreateChatCompletionStream(context.Background(), groq.ChatCompletionRequest{
				Model: groq.ModelLlama3170BVersatile,
				Messages: []groq.ChatCompletionMessage{
					{
						Role:    groq.ChatMessageRoleUser,
						Content: message,
					},
				},
				MaxTokens: 2000,
			})
			if err != nil {
				m.err = errMsg(err)
				return m, nil
			}
			newIx := len(m.messages) - 1
			currentCnt := ""
			for {
				response, err := re.Recv()
				if err != nil {
					m.err = errMsg(err)
					return m, nil
				}
				if response.Choices[0].FinishReason == groq.FinishReasonStop {
					break
				}
				currentCnt += response.Choices[0].Delta.Content
				m.messages[newIx] = m.senderStyle.Render("Groq: ") + currentCnt
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
			}
		}
	case tea.MouseAction:
		val := tea.MouseAction(msg)
		if val == 4 {
			m.viewport.LineUp(1)
			return m, nil
		} else if val == 5 {
			m.viewport.LineDown(1)
			return m, nil
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
