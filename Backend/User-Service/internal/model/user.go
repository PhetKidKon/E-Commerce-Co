package model

import "time"

// User map กับตาราง users (password เป็น hash, ไม่ส่งออก JSON)
type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      *string   `json:"username,omitempty"` // NULL ได้ (OAuth)
	Password      string    `json:"-"`                  // hash — ไม่ออก JSON เด็ดขาด
	FullName      string    `json:"fullName"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	ActiveFlag    bool      `json:"activeFlag"`
	Provider      string    `json:"provider"`
	EmailVerified bool      `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
