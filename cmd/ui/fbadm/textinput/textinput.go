package textinput

import (
	"context"
	"fmt"
	"strings"

	"owncli/cmd/csvtodb/schema"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	ti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/api/option"
)

type readCsvMsg struct{}

type readingCSVMsg struct {
	plans     map[int]schema.Plan
	workouts  map[int]schema.Workout
	exercises []schema.Exercise
}

type dataInsertedMsg struct{}

type Model struct {
	textInput ti.Model
	err       error
	output    string
	quit      bool
}

func InitialModel() Model {
	t := ti.New()
	t.Placeholder = "Enter your message here"
	t.Focus()
	t.CharLimit = 256
	t.Width = 30

	return Model{
		err:       nil,
		textInput: t,
	}
}

func getFBApp() (*firebase.App, error) {
	opt := option.WithCredentialsFile("/home/soratto/Documentos/apps/workout/credentials-workout.json")
	config := &firebase.Config{
		ProjectID: "home-workout-494f0",
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, nil
}

func sendMessage() string {
	ctx := context.Background()
	app, err := getFBApp()
	client, err := app.Messaging(ctx)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	message := &messaging.Message{
		Topic: "debug",
		// Token: "c3UMYWZgSHqsvOd8sEizVw:APA91bFQaf4Jy8M1T7BgSHOKxOo8zAyyXhpMjYjSxK8TPmv6ucnSxjnFxvuaa3HZa-5NSmXZKoN6NzGNdQRjJj42VabxG3LHbRNO9OmvIqBoj_4KC_dierDlQsckiVUW93k7mwyZmVfp",
		Notification: &messaging.Notification{
			Title: "Do exercises at home!",
			Body:  "Do exercises at home today, don't waste time on the couch!",
		},
	}

	res, err := client.Send(ctx, message)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("Successfully sent message: %v\n", res)
}

func (m Model) Init() tea.Cmd {
	return ti.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			m.output = sendMessage()
			return m, tea.Quit
		}

		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var str strings.Builder
	if m.output != "" {
		str.WriteString(m.output)
		return str.String()
	}
	str.WriteString(m.textInput.View())

	str.WriteString("\n")

	return str.String()
}
