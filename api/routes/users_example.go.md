# Example: How to Update routes/users.go to Use models and sql Packages

Replace your `routes/users.go` with this pattern:

```go
package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"shawty-ur/api/models"
	"shawty-ur/api/sql"
	"shawty-ur/app"

	"github.com/go-chi/chi/v5"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(r chi.Router, application *app.Application) {
	// Create the user store
	userStore := sql.NewUserStore(application.DbConnector)

	// Register routes
	r.Get("/users", listUsersHandler(userStore))
	r.Post("/register", createUserHandler(userStore))
	r.Post("/login", loginUserHandler(userStore))
	r.Get("/users/{id}", getUserHandler(userStore))
	r.Delete("/users/{id}", deleteUserHandler(userStore))
}

// listUsersHandler lists all users
func listUsersHandler(store *sql.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := store.List(r.Context())
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// createUserHandler creates a new user
func createUserHandler(store *sql.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload models.UserPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Create user model
		user := &models.User{
			Username: payload.Username,
			Email:    payload.Email,
		}

		// Hash password
		if err := user.HashPassword(payload.Password); err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Save to database
		if err := store.Create(r.Context(), user); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "user created",
			"user_id": user.ID,
		})
	}
}

// loginUserHandler handles user login
func loginUserHandler(store *sql.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload models.UserPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user from database
		user, err := store.GetByEmail(r.Context(), payload.Email)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// TODO: Verify password with bcrypt.CompareHashAndPassword
		// TODO: Generate JWT token

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "login successful",
			"user_id": user.ID,
		})
	}
}

// getUserHandler gets a user by ID
func getUserHandler(store *sql.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := store.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// deleteUserHandler deletes a user by ID
func deleteUserHandler(store *sql.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if err := store.Delete(r.Context(), id); err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
```

## Key Changes:

1. **Import models and sql packages**: `import "shawty-ur/api/models"` and `import "shawty-ur/api/sql"`
2. **Use `models.User` instead of `routes.User`**
3. **Create `UserStore` in route registration**: `userStore := sql.NewUserStore(application.DbConnector)`
4. **Pass `userStore` to handlers** instead of `application`
5. **Use store methods**: `store.Create()`, `store.GetByEmail()`, etc.
