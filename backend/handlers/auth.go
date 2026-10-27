package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"archival-vault/backend/database"
	"archival-vault/backend/middleware"
	"archival-vault/backend/models"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.JSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSONError(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if len(req.Username) < 3 {
		middleware.JSONError(w, "username must be at least 3 characters", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		middleware.JSONError(w, "password must be at least 6 characters", http.StatusBadRequest)
		return
	}
	if !strings.Contains(req.Email, "@") {
		middleware.JSONError(w, "invalid email address", http.StatusBadRequest)
		return
	}

	var existingID int64
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = ? OR email = ?", req.Username, req.Email).Scan(&existingID)
	if err == nil {
		middleware.JSONError(w, "username or email already registered", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		middleware.JSONError(w, "database query failure: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.JSONError(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	res, err := database.DB.Exec(
		"INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, 'OPERATOR')",
		req.Username, req.Email, string(hashedPassword),
	)
	if err != nil {
		middleware.JSONError(w, "failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userID, _ := res.LastInsertId()
	token, err := middleware.GenerateToken(userID, req.Username, "OPERATOR")
	if err != nil {
		middleware.JSONError(w, "failed to issue authentication token", http.StatusInternalServerError)
		return
	}

	middleware.JSONSuccess(w, http.StatusCreated, models.AuthResponse{
		Token: token,
		User: models.User{
			ID:       userID,
			Username: req.Username,
			Email:    req.Email,
			Role:     "OPERATOR",
		},
	}, "user registered successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.JSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSONError(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	identifier := strings.TrimSpace(req.Username)
	if identifier == "" || req.Password == "" {
		middleware.JSONError(w, "username/email and password required", http.StatusBadRequest)
		return
	}

	var user models.User
	var hashedPassword string
	row := database.DB.QueryRow(
		"SELECT id, username, email, password, role, created_at FROM users WHERE username = ? OR email = ?",
		identifier, identifier,
	)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &hashedPassword, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.JSONError(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		middleware.JSONError(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		middleware.JSONError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		middleware.JSONError(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	middleware.JSONSuccess(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	}, "login successful")
}

func Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		middleware.JSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	row := database.DB.QueryRow("SELECT id, username, email, role, created_at FROM users WHERE id = ?", claims.UserID)
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt); err != nil {
		middleware.JSONError(w, "user not found", http.StatusNotFound)
		return
	}

	middleware.JSONSuccess(w, http.StatusOK, user, "user profile loaded")
}
