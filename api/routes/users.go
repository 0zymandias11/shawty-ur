package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"shawty-ur/api/helper"
	"shawty-ur/api/models"
	"shawty-ur/api/utils"
	"shawty-ur/api/utils/db"
	"shawty-ur/app"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

// EXAMPLE FILE - Shows how to add new routes
// To use: rename to users.go and uncomment the registration in main.go

// RegisterUserRoutes registers all user-related routes
type password struct {
	text *string
	hash []byte
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email"`
	Password  password  `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (password *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Cannot Hash Password", "err: ", err)
		return err
	}
	password.hash = hash
	password.text = &text
	return nil
}

// func (user *User) hashedPassword(password string) error {
// 	if err := user.Password.Set(password); err != nil {
// 		slog.Error("Error in password Hash Set", "err: ", err)
// 		return err
// 	}
// 	return nil
// }

func RegisterUserRoutes(r chi.Router, application *app.Application) {
	// All user routes go here
	r.Get("/users", listUsersHandler(application))
	r.Post("/register", createUserHandler(application))
	r.Post("/login", loginUserHandler(application))
	r.Get("/users/{id}", getUserHandler(application))
	r.Delete("/users/{id}", deleteUserHandler(application))
}

// func (user *User) createUser(ctx context.Context, tx *sql.Tx) error {
// 	slog.Info("Creating User Query: ", "email: ", user.Email, "Password: ", user.Password.text)
// 	query := `INSERT INTO users(email, password, username)
// 			VALUES ($1, $2, $3) ON CONFLICT(email) DO NOTHING
// 			RETURNING id created_At, updated_at`

// 	err := tx.QueryRowContext(ctx, query,
// 		user.Email,
// 		user.Password.hash,
// 		user.Username).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

// 	if err == sql.ErrNoRows {
// 		//user already exists
// 		return nil
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func listUsersHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userStore := helper.NewUserStore(app.DbConnector)
		userList, err := userStore.ListUsers(r.Context())
		if err != nil {
			slog.Error("Error getting user list")
			utils.WriteJSON(w, http.StatusInternalServerError, "Internal server Error")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "success",
			"userList": userList,
		})
	}
}

func createUserHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		request := new(models.UserPayload)
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			slog.Error("Error Decoding user request body !!", "err: ", err)
			utils.WriteJSON(w, http.StatusBadRequest, "BAD Login Request")
			return
		}
		user := new(models.User)
		if err := user.HashPassword(request.Password); err != nil {
			slog.Error("Error in hashing User Password ", "err: ", err)
			utils.WriteJSON(w, http.StatusInternalServerError, "Failed to process password")
			return
		}
		user.Email = request.Email
		if request.Username == "" {
			user.Username = request.Email
		} else {
			user.Username = request.Username
		}
		user.Username = request.Username
		slog.Info("Creating User: ", "email: ", user.Email, "Password: ", user.Password.Text)
		txErr := db.WithTx(app.DbConnector, req.Context(), func(tx *sql.Tx) error {
			userStore := helper.NewUserStore(app.DbConnector)

			if err := userStore.CreateUser(req.Context(), tx, user); err != nil {
				slog.Error("Error running create User query !!", "err", err)
				return err // Return error to rollback transaction
			}
			return nil
		})

		if txErr != nil {
			slog.Error("Error in create user tx!!! ", "err", txErr)
			utils.WriteJSON(w, http.StatusInternalServerError, "Error creating new User!!!")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "user created",
		})
	}
}

func loginUserHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		request := new(UserPayload)
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, "Invalid User Creation Request Body")
		}

		user := new(models.User)
		user.Email = request.Email
		user.Username = request.Username
		user.Password.Text = request.Password

		txErr := db.WithTx(app.DbConnector, req.Context(), func(tx *sql.Tx) error {
			userStore := helper.NewUserStore(app.DbConnector)

			err := userStore.GetUserByUsernameOrEmail(req.Context(), tx, user)
			if err != nil {
				log.Printf("ERROR:: Cannot Find user with username: %s, email: %s", user.Username, user.Email)
				utils.WriteJSON(w, http.StatusBadRequest, "Invalid Username/Email")
			}

			if err := user.CompareHash(); err != nil {
				slog.Error("Cannot login Invalid Password", "err: ", err)
				utils.WriteJSON(w, http.StatusForbidden, "Invalid Username/Password")
				return nil
			}
			return nil
		})

		if txErr != nil {
			slog.Error("Error in create user tx!!! ", "err", txErr)
			utils.WriteJSON(w, http.StatusInternalServerError, "Error creating new User!!!")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "success",
			"ID":       user.ID,
			"username": user.Username,
		})
	}
}

func getUserHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get user: " + id))
	}
}

func deleteUserHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		fmt.Println("id: ", id)
		w.WriteHeader(http.StatusNoContent)
	}
}
