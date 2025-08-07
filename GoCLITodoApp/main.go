package main

import (
	"fmt"
	"os"
	"log"
	"context"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type todo struct {
	ID			primitive.ObjectID		`bson:"_id,omitempty"`
	Completed	bool					`bson:"completed"`
	Text		string					`bson:"text"`
}

type model struct {
	Cursor		int
	Choices 	[]string
	Mode		mode
	Todos		[]todo
	InputBuf	string
}

type mode int

const (
	menuMode mode = iota
	newTodoMode
	viewTodosMode
)

const clearScreen = "\033[H\033[2J"

var todoCollection *mongo.Collection

func initialModel() model {
	todos := getTodos()
	log.Println(len(todos))
	return model{
		Choices: []string{"New todo", "View todos", "Exit"},
		Todos: todos,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
			
		case "ctrl+c", "q":
			if m.Mode == menuMode {
				setTodos(m.Todos)
				return m, tea.Quit
			} else if m.Mode == newTodoMode {
				m.InputBuf += string(msg.Runes)
			}

		case "up", "k":
			if m.Cursor > 0 && m.Mode != newTodoMode {
				m.Cursor--
			} else {
				m.InputBuf += string(msg.Runes)
			}

		case "down", "j":
			if m.Cursor < len(m.Choices)-1 && m.Mode == menuMode {
				m.Cursor++
			} else if m.Cursor < len(m.Todos)-1 && m.Mode == viewTodosMode {
				m.Cursor++
			} else if m.Mode == newTodoMode {
				m.InputBuf += string(msg.Runes)
			}

		case "d":
			if m.Mode == viewTodosMode && len(m.Todos) > 1{
				deleteTodo(m.Todos[m.Cursor].ID)
				m.Todos = append(m.Todos[:m.Cursor], m.Todos[m.Cursor+1:]...)
			} else if m.Mode == newTodoMode {
				m.InputBuf += string(msg.Runes)
			}


		case "esc":
			m.Mode = menuMode

		case " ":
			if m.Mode == viewTodosMode && len(m.Todos) > 0 {
				m.Todos[m.Cursor].Completed = !m.Todos[m.Cursor].Completed
			} else if m.Mode == newTodoMode {
				m.InputBuf += string(msg.Runes)
			}

		case "enter":
			switch m.Mode {

				case menuMode:
					switch m.Choices[m.Cursor] {
				
					case "Exit":
						setTodos(m.Todos)
						return m, tea.Quit

					case "New todo":
						m.Mode = newTodoMode
						m.Cursor = 0

					case "View todos":
						m.Mode = viewTodosMode
						m.Cursor = 0
	
					}
					
				case newTodoMode:
					if len(m.InputBuf) > 0 {
						m.Todos = append(m.Todos, todo{
							ID:			primitive.NewObjectID(),
							Completed:	false,
							Text:		m.InputBuf,
						})

						m.InputBuf = ""
						m.Mode = menuMode
						m.Cursor = 0
					}
			} 

		case "backspace":
			if m.Mode == newTodoMode {
				if len(m.InputBuf) > 0 {
					m.InputBuf = m.InputBuf[:len(m.InputBuf)-1]
				}
			}

		default:
			if m.Mode == newTodoMode {
				if len(msg.Runes) > 0 {
					m.InputBuf += string(msg.Runes)
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := clearScreen

	switch m.Mode {

	case menuMode:
		s += "Hello, welcome to CLI todo app!\n\n"

		for i, choice := range m.Choices {
			cursor := " "
			if m.Cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
	
		s += "\nPress q to quit. \n"

	case newTodoMode:
		s += "New todo: " 

		s += m.InputBuf + "_"
		
		s += "\nPress ESC to go back."
		
	case viewTodosMode:
		s += "Your todos: \n"

		if len(m.Todos) == 0 {
			s += "-- No todos yet. -- \n"
		} else {
			for i, t := range m.Todos {
				cursor := " "
				if m.Cursor == i {
					cursor = ">"
				}

				check := "[ ]"
				if t.Completed {
					check = "[x]"
				}
				s += fmt.Sprintf("%s %s %s\n", cursor, check, t.Text)
			}

			s += "\nPress space to toggle completion." 
		}

		s += "\nPress ESC to go back."
	}


	return s
}

func getTodos() ([]todo) {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	cursor, err := todoCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var todos []todo
	for cursor.Next(ctx) {
		var t todo
		if err := cursor.Decode(&t); err != nil {
			return nil
		}
		todos = append(todos, t)
	}

	if err:= cursor.Err(); err != nil {
		return nil
	}

	return todos
}

func setTodos(todos []todo) {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	for _, t := range todos {
		filter := bson.M{"_id": t.ID}
		update := bson.M{"$set": t}
		opts := options.Update().SetUpsert(true)
		_, err := todoCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func deleteTodo(id primitive.ObjectID) {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	filter := bson.M{"_id": id}

	res, err := todoCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	if res.DeletedCount == 0 {
		log.Println("No document found wiht id:", id)
	}
}


func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("MONGO_USERNAME")
	dbPass := os.Getenv("MONGO_PASSWORD")

	f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
    	log.Fatal(err)
	}
	log.SetOutput(f)

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.uv4lr.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", dbUser, dbPass)))
	if err != nil {
		log.Fatal(err)
	} 
	defer client.Disconnect(ctx)
	
	todoCollection = client.Database("temp01").Collection("todos")

		

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}
