# Final Project PBKK

## Overview

This web application, built with Go, allows users to create, read, update, and delete code snippets. It features user authentication, secure password storage, and a responsive interface designed with Tailwind CSS.

## Features

- User Registration and Authentication 
- Create, Read, Update, and Delete (CRUD) Operations for Code Snippets
- Secure Password Hashing with bcrypt
- Session Management with JWT Tokens
- Responsive Design Using Tailwind CSS


## Prerequisites

- Go (version 1.23 or later)
- Git
- SQLite3

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/azarelga/final-project-pbkk.git
cd final-project-pbkk
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment Variables

Create a `.env` file in the project root with the following content:

```
PORT=8080
SESSION_KEY=your_very_long_and_random_secret_key_here
DATABASE_PATH=./snippets.db
```

Note: Replace `your_very_long_and_random_secret_key_here` with a secure, random string of at least 32 characters. You can generate one using:

```bash
openssl rand -base64 32
```

## Running the Application

### Development Mode

```bash
go run main.go
```

### Build and Run

```bash
go build -o code-snippets
./code-snippets
```

## Project Structure

```
.
├── main.go                # Main application entry point
├── .env                   # Environment configuration
├── go.mod                 # Go module dependencies
├── go.sum                 # Go module checksums
├── handlers/              # HTTP request handlers
│   ├── auth.go
│   └── snippets.go
├── repositories/          # Database access layers
│   ├── user.go
│   └── snippets.go
├── services/              # Business logic
│   ├── user.go
│   └── snippets.go
├── middleware/            # Middleware functions
│   └── checkAuth.go
├── database/              # Database initialization and migrations
│   ├── db.go
│   ├── migrate.go
│   └── loadenvs.go
├── templates/             # HTML templates
│   ├── header.html
│   ├── footer.html
│   ├── home.html
│   ├── login.html
│   ├── register.html
│   ├── list.html
│   ├── mylist.html
│   ├── create.html
│   ├── edit.html
│   └── viewsnippet.html
├── .gitignore             # Git ignore file
├── README.md              # Project documentation
└── snippets.db            # SQLite database (generated at runtime)
```
## CRUD Implementation
The application implements CRUD (Create, Read, Update, Delete) operations for code snippets:

- **Create**: Users can create new code snippets using the `CreateSnippet` handler in `snippets.go`. The HTML form for creating snippets is in `create.html`.

- **Read**: Users can read existing snippets through several handlers in `snippets.go`:
  - `GetSnippetsByLanguage` retrieves snippets filtered by language.
  - `GetSnippetsByUserID` lists snippets created by the current user.
  - `GetSnippetByID` displays detailed information about a specific snippet.
- **Update**: Users can update their snippets using the `UpdateSnippet` handler in `snippets.go`. The edit form is provided in `edit.html`.

- **Delete**: Users can delete their snippets using the `DeleteSnippet` handler in `snippets.go`.

These handlers interact with the `SnippetService` in `snippets.go`, which uses the `SnippetRepository` in `snippets.go` for database operations. Authentication and authorization are managed by middleware in `checkAuth.go` to ensure that only authorized users can perform these actions.

## Security Considerations

- Password Hashing: User passwords are securely hashed using bcrypt.
- Session Management: Authentication is managed with JWT tokens stored in cookies.
- Protected Routes: Middleware ensures that certain routes are only accessible to authenticated users.

## Customization

### Changing the Database
By default, the application uses SQLite. To switch to a different database (e.g., MySQL or PostgreSQL), update the database configuration in db.go and adjust the GORM dialect accordingly.

### Styling
The application uses Tailwind CSS for styling. You can customize the styles by modifying the HTML templates in the templates directory.

## Demo Project
https://youtu.be/3ma7tgqgtKg