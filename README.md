# Final Project PBKK

## Overview

This is a web application built with Go that allows users to create, read, update, and delete code snippets. The application features user authentication, secure password storage, and a clean, responsive interface.

## Features

- User Registration and Authentication
- Create, Read, Update, and Delete (CRUD) Code Snippets
- Secure Password Hashing
- Session-based Authentication
- Responsive Design with Tailwind CSS

## Prerequisites

- Go (version 1.16 or later)
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
code-snippet-sharing/
│
├── main.go               # Main application logic
├── .env                  # Environment configuration
├── go.mod                # Go module dependencies
│
├── templates/            # HTML templates
│   ├── base.html
│   ├── home.html
│   ├── login.html
│   ├── register.html
│   ├── list.html
│   ├── create.html
│   ├── view.html
│   └── edit.html
│
└── snippets.db           # SQLite database (auto-created)
```

## Security Considerations

- Passwords are hashed using bcrypt
- Sessions are managed securely
- Authentication middleware protects routes
- SQLite database is used for lightweight, file-based storage

## Customization

### Changing the Database

- The application uses SQLite by default
- To use a different database, modify the database initialization in `main.go`

### Styling

- Tailwind CSS is used for styling
- Modify templates to customize the look and feel
