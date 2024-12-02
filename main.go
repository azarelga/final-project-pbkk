package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user model
type User struct {
	ID       int
	Username string
	Password string
}

// CodeSnippet represents the structure of a code snippet
type CodeSnippet struct {
	ID          int
	Title       string
	Language    string
	Code        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Database connection
var (
	db           *sql.DB
	sessionStore *sessions.CookieStore
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using defaults")
	}

	// Initialize database
	initDatabase()
	defer db.Close()

	// Setup router
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Authentication routes
	r.HandleFunc("/register", registerHandler).Methods("GET", "POST")
	r.HandleFunc("/login", loginHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")

	// Routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/snippets", listSnippetsHandler).Methods("GET")
	r.HandleFunc("/snippets/new", createSnippetFormHandler).Methods("GET")
	r.HandleFunc("/snippets/new", createSnippetHandler).Methods("POST")
	r.HandleFunc("/snippets/{id:[0-9]+}", viewSnippetHandler).Methods("GET")
	r.HandleFunc("/snippets/{id:[0-9]+}/edit", editSnippetFormHandler).Methods("GET")
	r.HandleFunc("/snippets/{id:[0-9]+}/edit", updateSnippetHandler).Methods("POST")
	r.HandleFunc("/snippets/{id:[0-9]+}/delete", deleteSnippetHandler).Methods("POST")

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Database Initialization
func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./snippets.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create snippets table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS snippets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			language TEXT NOT NULL,
			code TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

// Handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/home.html"))
	tmpl.Execute(w, nil)
}

func listSnippetsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, language, description, created_at FROM snippets ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var snippets []CodeSnippet
	for rows.Next() {
		var s CodeSnippet
		err := rows.Scan(&s.ID, &s.Title, &s.Language, &s.Description, &s.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		snippets = append(snippets, s)
	}

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/list.html"))
	tmpl.Execute(w, snippets)
}

func createSnippetFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/create.html"))
	tmpl.Execute(w, nil)
}

func createSnippetHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.Form.Get("title")
	language := r.Form.Get("language")
	code := r.Form.Get("code")
	description := r.Form.Get("description")

	result, err := db.Exec(
		"INSERT INTO snippets (title, language, code, description) VALUES (?, ?, ?, ?)",
		title, language, code, description,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}

func viewSnippetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var snippet CodeSnippet
	err := db.QueryRow(
		"SELECT id, title, language, code, description, created_at FROM snippets WHERE id = ?",
		id,
	).Scan(
		&snippet.ID, &snippet.Title, &snippet.Language,
		&snippet.Code, &snippet.Description, &snippet.CreatedAt,
	)
	if err != nil {
		http.Error(w, "Snippet not found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/view.html"))
	tmpl.Execute(w, snippet)
}

func editSnippetFormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var snippet CodeSnippet
	err := db.QueryRow(
		"SELECT id, title, language, code, description FROM snippets WHERE id = ?",
		id,
	).Scan(
		&snippet.ID, &snippet.Title, &snippet.Language,
		&snippet.Code, &snippet.Description,
	)
	if err != nil {
		http.Error(w, "Snippet not found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/edit.html"))
	tmpl.Execute(w, snippet)
}

func updateSnippetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.Form.Get("title")
	language := r.Form.Get("language")
	code := r.Form.Get("code")
	description := r.Form.Get("description")

	_, err = db.Exec(
		"UPDATE snippets SET title = ?, language = ?, code = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		title, language, code, description, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}

func deleteSnippetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	_, err := db.Exec("DELETE FROM snippets WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/snippets", http.StatusSeeOther)
}

// Authentication Middleware
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessionStore.Get(r, "session")

		// Check if user is authenticated
		userID, ok := session.Values["user_id"].(int)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Add user ID to request context (optional, but can be useful)
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Registration Handler
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/register.html"))
		tmpl.Execute(w, nil)
		return
	}

	// Handle POST request
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// Validate input
	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Redirect to login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login Handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}

	// Handle POST request
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// Find user
	var user User
	err = db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&user.ID, &user.Password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create session
	session, _ := sessionStore.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	// Redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout Handler
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, "session")

	// Clear the session
	session.Values["user_id"] = nil
	session.Options.MaxAge = -1
	session.Save(r, w)

	// Redirect to login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
