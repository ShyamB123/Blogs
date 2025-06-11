# Blogs

A simple blog platform built with Go (backend) and HTML/CSS/JavaScript (frontend).

## Features

- View all blogs
- View blog details
- Add, edit, and remove blogs (requires authentication)
- Simple user authentication (username/password)
- Data stored in JSON files

## Project Structure

```
Blogs/
├── backend/
│   ├── main.go
│   ├── blogs.json
│   ├── users.json
│   └── .env
├── frontend/
│   ├── index.html
│   ├── blog.html
│   └── style.css
└── README.md
```

## Getting Started

### Backend

1. **Install Go** (if not already installed):  
   https://golang.org/dl/

2. **Install dependencies**:  
   ```
   go get github.com/joho/godotenv
   ```

3. **Set up `.env` file** in `backend/`:
   ```
   AUTH=false
   ```

4. **Run the server**:
   ```
   cd backend
   go run main.go
   ```

   The server runs at [http://localhost:8080](http://localhost:8080).

### Frontend

Open `frontend/index.html` in your browser.

## API Endpoints

- `GET /blogs` — List all blogs
- `POST /add` — Add a new blog (requires `AUTH=true`)
- `POST /edit` — Edit a blog (requires `AUTH=true`)
- `POST /remove` — Remove a blog (requires `AUTH=true`)
- `POST /login` — Log in with username and password
- `POST /logout` — Log out

## Authentication

- By default, add/edit/remove endpoints are **disabled** (`AUTH=false`).
- After a successful login, `AUTH` is set to `true` and those endpoints become available.
