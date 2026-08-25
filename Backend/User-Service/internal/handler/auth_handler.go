package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kidkon/ecommerce/common/response"
	"github.com/kidkon/ecommerce/user-service/internal/authmw"
	"github.com/kidkon/ecommerce/user-service/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Fail(w, response.BadRequest("invalid json body"))
		return
	}
	u, err := h.svc.Register(r.Context(), service.RegisterInput{
		Email: body.Email, Username: body.Username, Password: body.Password, FullName: body.FullName,
	})
	if err != nil {
		response.Fail(w, err)
		return
	}
	response.Created(w, u)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Fail(w, response.BadRequest("invalid json body"))
		return
	}
	res, err := h.svc.Login(r.Context(), body.Login, body.Password)
	if err != nil {
		response.Fail(w, err)
		return
	}
	response.OK(w, res)
}

// Me: ข้อมูลเบื้องต้นจาก "token" (เร็ว ไม่แตะ DB) — snapshot ตอน login
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	c := authmw.ClaimsFrom(r.Context())
	response.OK(w, map[string]any{
		"userId":        c.UserID,
		"username":      c.Username,
		"fullName":      c.FullName,
		"email":         c.Email,
		"role":          c.Role,
		"provider":      c.Provider,
		"emailVerified": c.EmailVerified,
	})
}

// Profile: ข้อมูลละเอียดจาก "DB" สด (ด้วย userId จาก token)
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	c := authmw.ClaimsFrom(r.Context())
	u, err := h.svc.GetProfile(r.Context(), c.UserID)
	if err != nil {
		response.Fail(w, err)
		return
	}
	response.OK(w, u) // full model.User สด (password ไม่ออกเพราะ json:"-")
}
