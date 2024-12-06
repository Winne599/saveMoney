package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
)

type User struct {
	UserID           int        `json:"userId"`           // 用户ID
	Username         string     `json:"username"`         // 用户名
	PasswordHash     string     `json:"passwordHash"`     // 密码哈希
	Email            string     `json:"email"`            // 邮箱
	PhoneNumber      string     `json:"phoneNumber"`      // 手机号
	RegistrationTime time.Time  `json:"registrationTime"` // 注册时间
	LastLoginTime    *time.Time `json:"lastLoginTime"`    // 最后登录时间
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO `users` (`username`, `password_hash`, `email`, `phone_number`) VALUES (?, ?, ?, ?)",
		user.Username, user.PasswordHash, user.Email, user.PhoneNumber)
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok && mysqlErr.Number == 1062 {
			http.Error(w, "Username already exists", http.StatusBadRequest)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintln(w, "User registered successfully")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var dbUser User
	err = db.QueryRow("SELECT user_id, password_hash FROM users WHERE username = ?", user.Username).Scan(&dbUser.UserID, &dbUser.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if user.PasswordHash != dbUser.PasswordHash {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// TODO: Generate and store session token in Redis

	fmt.Fprintln(w, "Login successful")
}
