package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"

	"navmate-backend/config"
	"navmate-backend/internal/models"
	"navmate-backend/pkg/hash"
	"navmate-backend/pkg/jwtauth"
)

type GoogleHandler struct {
	db    *gorm.DB
	jwt   *jwtauth.Service
	oauth *oauth2.Config
}

func NewGoogleHandler(db *gorm.DB, jwt *jwtauth.Service, cfg *config.Config) *GoogleHandler {
	return &GoogleHandler{
		db:  db,
		jwt: jwt,
		oauth: &oauth2.Config{
			ClientID:     cfg.Google.ClientID,
			ClientSecret: cfg.Google.ClientSecret,
			RedirectURL:  cfg.Google.RedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (h *GoogleHandler) Login(c *gin.Context) {
	state := randomState(32)
	// เก็บ state ไว้ใน cookie (dev: Secure=false; prod ควร true + SameSite)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(5 * time.Minute),
	})
	url := h.oauth.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func (h *GoogleHandler) Callback(c *gin.Context) {
	// verify state
	stateQ := c.Query("state")
	ck, _ := c.Request.Cookie("oauthstate")
	if ck == nil || ck.Value == "" || ck.Value != stateQ {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	ctx := context.Background()
	tok, err := h.oauth.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "code exchange failed"})
		return
	}

	// ดึงข้อมูลผู้ใช้จาก Google
	gu, err := fetchGoogleUser(tok.AccessToken)
	if err != nil || gu.Email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "get userinfo failed"})
		return
	}

	// หา/สร้างผู้ใช้ในระบบเรา
	email := strings.ToLower(gu.Email)
	var u models.User
	if err := h.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// สร้าง user ใหม่จาก Google
			dummyHash, _ := hash.HashPassword(randomState(24))
			u = models.User{
				Email:        email,
				PasswordHash: dummyHash,
				Provider:     "google",
			}
			if gu.ID != "" {
				gid := gu.ID
				u.GoogleID = &gid
			}
			if err := h.db.Create(&u).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "create user failed"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
	} else {
		// มีอยู่แล้ว → อัปเดต provider/google_id ถ้าไม่มี
		changed := false
		if u.Provider != "google" {
			u.Provider = "google"
			changed = true
		}
		if u.GoogleID == nil && gu.ID != "" {
			gid := gu.ID
			u.GoogleID = &gid
			changed = true
		}
		if changed {
			_ = h.db.Save(&u).Error
		}
	}

	// ออก JWT ของเรา
	token, err := h.jwt.GenerateToken(u.ID, u.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	// ส่งกลับเป็น JSON (ถ้าอยาก redirect ไป FE พร้อม token ก็ทำได้)
	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"email":    u.Email,
		"provider": u.Provider,
	})
}

type googleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func fetchGoogleUser(accessToken string) (*googleUser, error) {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var gu googleUser
	if err := json.NewDecoder(resp.Body).Decode(&gu); err != nil {
		return nil, err
	}
	return &gu, nil
}

func randomState(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
