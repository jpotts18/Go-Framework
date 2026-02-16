package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"golaravel/app"
	"golaravel/app/http/middleware"
	"golaravel/app/http/request"
	"golaravel/app/http/response"
	"golaravel/app/routing"
	"golaravel/app/validation"
)

// Todo represents a single todo item.
type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// TodoStore is a simple in-memory store for todos.
type TodoStore struct {
	mu     sync.RWMutex
	todos  []Todo
	nextID int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{nextID: 1}
}

func (s *TodoStore) All() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Todo, len(s.todos))
	copy(result, s.todos)
	return result
}

func (s *TodoStore) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	todo := Todo{ID: s.nextID, Title: title}
	s.nextID++
	s.todos = append(s.todos, todo)
	return todo
}

func (s *TodoStore) Find(id int) (Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.todos {
		if t.ID == id {
			return t, true
		}
	}
	return Todo{}, false
}

func (s *TodoStore) Toggle(id int) (Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.todos {
		if t.ID == id {
			s.todos[i].Completed = !s.todos[i].Completed
			return s.todos[i], true
		}
	}
	return Todo{}, false
}

func (s *TodoStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.todos {
		if t.ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return true
		}
	}
	return false
}

func main() {
	application := app.New(".")
	application.SetViewPath("./resources/views")

	store := NewTodoStore()

	// Seed a couple of example todos
	store.Create("Learn GoLaravel")
	store.Create("Build something cool")

	router := application.GetRouter()

	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// ── HTML Routes ─────────────────────────────────────────────
	router.Get("/", indexHandler(application, store))
	router.Post("/todos", createHandler(application, store))
	router.Post("/todos/{id}/toggle", toggleHandler(store))
	router.Post("/todos/{id}/delete", deleteHandler(store))

	// ── JSON API ────────────────────────────────────────────────
	router.Group("/api", func(g *routing.RouteGroup) {
		g.Use(middleware.JSON())

		g.Get("/todos", apiListHandler(store))
		g.Post("/todos", apiCreateHandler(store))
		g.Put("/todos/{id}/toggle", apiToggleHandler(store))
		g.Delete("/todos/{id}", apiDeleteHandler(store))
	})

	app.PrintBanner()
	fmt.Println("  Todo app running at http://localhost:5000")
	fmt.Println()
	application.Run(":5000")
}

// ── HTML Handlers ───────────────────────────────────────────────

func indexHandler(application *app.Application, store *TodoStore) routing.HandlerFunc {
	return func(req *request.Request, res *response.Response) error {
		todos := store.All()

		html, err := application.GetView().Render("index", map[string]interface{}{
			"title": "Todo App",
			"todos": todos,
			"count": len(todos),
		})
		if err != nil {
			return res.ServerError(err.Error())
		}

		return res.HTML(html)
	}
}

func createHandler(application *app.Application, store *TodoStore) routing.HandlerFunc {
	return func(req *request.Request, res *response.Response) error {
		title := req.FormValue("title")

		v := validation.Make(
			map[string]interface{}{"title": title},
			map[string][]string{"title": {"required", "string", "min:1", "max:200"}},
		)
		if v.Fails() {
			// Re-render the page with errors
			todos := store.All()
			html, err := application.GetView().Render("index", map[string]interface{}{
				"title":  "Todo App",
				"todos":  todos,
				"count":  len(todos),
				"errors": v.All(),
			})
			if err != nil {
				return res.ServerError(err.Error())
			}
			return res.Status(http.StatusUnprocessableEntity).HTML(html)
		}

		store.Create(title)

		return res.Redirect("/", http.StatusSeeOther)
	}
}

func toggleHandler(store *TodoStore) routing.HandlerFunc {
	return func(req *request.Request, res *response.Response) error {
		id, err := strconv.Atoi(req.Param("id"))
		if err != nil {
			return res.BadRequest("Invalid todo ID")
		}

		if _, ok := store.Toggle(id); !ok {
			return res.NotFound("Todo not found")
		}

		return res.Redirect("/", http.StatusSeeOther)
	}
}

func deleteHandler(store *TodoStore) routing.HandlerFunc {
	return func(req *request.Request, res *response.Response) error {
		id, err := strconv.Atoi(req.Param("id"))
		if err != nil {
			return res.BadRequest("Invalid todo ID")
		}

		if !store.Delete(id) {
			return res.NotFound("Todo not found")
		}

		return res.Redirect("/", http.StatusSeeOther)
	}
}

// ── JSON API Handlers ───────────────────────────────────────────

func apiListHandler(store *TodoStore) routing.HandlerFunc {
	return func(req *request.Request, res *response.Response) error {
		return res.JSON(store.All())
	}
}

func apiCreateHandler(store *TodoStore) routing.HandlerFunc {
	return func(req *request.Request, res *response.Response) error {
		v := validation.Make(req.All(), map[string][]string{
			"title": {"required", "string", "min:1", "max:200"},
		})
		if v.Fails() {
			return res.ValidationError(v.Errors())
		}

		title, _ := req.All()["title"].(string)
		todo := store.Create(title)
		return res.Status(http.StatusCreated).JSON(todo)
	}
}

func apiToggleHandler(store *TodoStore) routing.HandlerFunc {
	return func(req *request.Request, res *response.Response) error {
		id, err := strconv.Atoi(req.Param("id"))
		if err != nil {
			return res.BadRequest("Invalid todo ID")
		}

		todo, ok := store.Toggle(id)
		if !ok {
			return res.NotFound("Todo not found")
		}

		return res.JSON(todo)
	}
}

func apiDeleteHandler(store *TodoStore) routing.HandlerFunc {
	return func(req *request.Request, res *response.Response) error {
		id, err := strconv.Atoi(req.Param("id"))
		if err != nil {
			return res.BadRequest("Invalid todo ID")
		}

		if !store.Delete(id) {
			return res.NotFound("Todo not found")
		}

		return res.NoContent()
	}
}
