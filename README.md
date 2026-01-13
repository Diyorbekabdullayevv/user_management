# MongoDB User Management Project

A modern full-stack application for user registration and filtering by age range using Go backend and responsive frontend.

## Project Structure

```
MongoDB-Project/
├── backend/                 # Go backend code
│   ├── cmd/
│   │   └── main.go         # Server entry point
│   ├── db/
│   │   └── db.go           # MongoDB connection
│   ├── internal/
│   │   ├── docktor.go      # API endpoint
│   │   └── user.go         # User operations (CRUD, filter)
│   ├── models/
│   │   └── models.go       # User data model
│   ├── go.mod              # Go module definition
│   └── go.sum              # Go dependencies checksum
│
├── frontend/               # Frontend files
│   ├── html/
│   │   └── index.html      # Main HTML template
│   ├── css/
│   │   └── style.css       # Styling
│   └── js/
│       └── script.js       # JavaScript logic
│
├── cmd/                    # Legacy (old structure)
├── db/                     # Legacy (old structure)
├── internal/               # Legacy (old structure)
├── models/                 # Legacy (old structure)
├── go.mod                  # Legacy (old structure)
│
├── commands.txt            # Command reference
├── text.txt                # Notes
└── README.md               # This file
```

## Setup & Running

### Prerequisites
- Go 1.23+
- MongoDB running on `mongodb://localhost:27017`
- Docker (for running MongoDB)

### Start MongoDB
```bash
docker start mongodb
```

### Start Backend Server
```bash
cd backend
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## Features

### User Registration
- Register new users with Name, Last Name, and Age
- Form validation
- Automatic redirect to search page after registration
- Success/error notifications

### User Filtering
- Filter users by age range (min_age to max_age)
- Display results in beautiful cards
- Shows user count and details
- No results message when no users match

### API Endpoints

#### User Registration
```
POST /api/users
Content-Type: application/json

{
  "name": "John",
  "surname": "Doe",
  "age": 30
}
```

#### Filter Users by Age
```
GET /api/users/filter?min_age=20&max_age=40
```

#### Other Available Endpoints
- `POST /users` - Create single user
- `GET /users/:id` - Get user by ID
- `PUT /users/update/:id` - Update user
- `DELETE /users/delete/:id` - Delete user
- `POST /users/bulk` - Bulk create users
- `GET /users/all?status=claim&greater-than=X` - Get users with conditions

## Frontend Structure

### HTML (`frontend/html/index.html`)
- Semantic HTML5
- Two-tab interface (Register & Filter)
- Form elements with validation
- Results container for user display

### CSS (`frontend/css/style.css`)
- Professional gradient design
- Responsive layout
- Tab-based navigation with animations
- Card-based user display
- Success/error message styling
- Smooth transitions and hover effects

### JavaScript (`frontend/js/script.js`)
- Tab switching logic
- User registration form handling
- Age range filtering
- API communication via Fetch
- Dynamic user card rendering
- Form validation
- Loading and error states

## Development Notes

The project was reorganized to separate concerns:
- **Backend**: All Go code (database, models, API routes)
- **Frontend**: All HTML, CSS, and JavaScript (user interface)

The backend serves the frontend files via Gin static routes:
- `/css/*` → serves CSS files
- `/js/*` → serves JavaScript files
- `/` → serves index.html

## Browser Compatibility

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- Modern mobile browsers

## Notes

- The original files in `cmd/`, `db/`, `internal/`, `models/`, and `go.mod` are legacy and can be deleted
- All active development should use files in the `backend/` folder
- CSS and JS are now separate from HTML for better maintainability
