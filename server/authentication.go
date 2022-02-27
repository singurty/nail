package server

import (
	"net/http"
	"time"
	"context"
	"crypto/rand"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	pgp "github.com/ProtonMail/gopenpgp/v2/crypto"

	"github.com/singurty/nail/user"
	"github.com/singurty/nail/db"
	"github.com/singurty/nail/frontend"
)

func loginHandler(c *gin.Context) {
	ctx := c.MustGet("session").(context.Context)

	if c.Request.Method == "GET" {
		// Redirect to root if already logged in
		if sessionManager.Exists(ctx, "user_id") {
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}

	username, result := c.GetPostForm("username")
	if !result || username == "" {
		c.String(http.StatusOK, "Username should not be empty")
		return
	}
	password, result := c.GetPostForm("password")
	if !result || password == "" {
		c.String(http.StatusOK, "Password should not be empty")
		return
	}
	id, twoFactor, err := user.Login(username, password)
	if err != nil {
		c.String(http.StatusOK, "Username or password incorrect")
		return
	}

	// Add login to session
	sessionManager.Put(ctx, "user_id", id)
	
	// Two-factor authentication
	if twoFactor {
		pgpKeyRow := db.DBpool.QueryRow(context.Background(), "SELECT pgp_key FROM users WHERE id = $1;", id)
		var pgpKeyString string
		err := pgpKeyRow.Scan(&pgpKeyString)
		if err != nil {
			log.Error(err)
			c.String(http.StatusInternalServerError, "An error occurred while processing two-factor authentication")
			return
		}
		otp := generateOtp()
		pgpKey, err := pgp.NewKeyFromArmored(pgpKeyString)
		if err != nil {
			log.Error(err)
			c.String(http.StatusInternalServerError, "An error occurred while processing two-factor authentication")
			return
		}
		pgpKeyRing, err := pgp.NewKeyRing(pgpKey)
		if err != nil {
			log.Error(err)
			c.String(http.StatusInternalServerError, "An error occurred while processing two-factor authentication")
			return
		}
		pgpMessage, err := pgpKeyRing.Encrypt(pgp.NewPlainMessage([]byte(otp)), nil)
		if err != nil {
			log.Error(err)
			c.String(http.StatusInternalServerError, "An error occurred while processing two-factor authentication")
			return
		}
		pgpArmored, err := pgpMessage.GetArmored()
		if err != nil {
			log.Error(err)
			c.String(http.StatusInternalServerError, "An error occurred while processing two-factor authentication")
			return
		}
		message := frontend.TwoFactorPage{Message:pgpArmored}
		c.HTML(http.StatusOK, "twofactor.html", message)

		sessionManager.Put(ctx, "authorized", false)
		sessionToken, _, err := sessionManager.Commit(ctx)
		if err != nil {
			log.Error(err)
			c.String(http.StatusInternalServerError, "An error occured")
			return
		}
		sessionManager.WriteSessionCookie(ctx, c.Writer, sessionToken, time.Now())
	} else {
		sessionManager.Put(ctx, "authorized", true)
		sessionToken, _, err := sessionManager.Commit(ctx)
		if err != nil {
			log.Error(err)
			c.String(http.StatusInternalServerError, "An error occured")
			return
		}
		sessionManager.WriteSessionCookie(ctx, c.Writer, sessionToken, time.Now())
		c.Redirect(http.StatusSeeOther, "/")
	}
}

func generateOtp() string {
	otpChars := "0123456789"
	buffer := make([]byte, 6)
    _, err := rand.Read(buffer)
    if err != nil {
        log.Fatal(err)
    }

    otpCharsLength := len(otpChars)
    for i := 0; i < 6; i++ {
        buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
    }

    return string(buffer)
}

func registerHandler(c *gin.Context) {
	ctx := c.MustGet("session").(context.Context)

	if c.Request.Method == "GET" {
		// Redirect to root if already logged in
		if sessionManager.Exists(ctx, "user_id") {
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.HTML(http.StatusOK, "register.html", nil)
		return
	}

	username, result := c.GetPostForm("username")
	if !result || username == "" {
		c.String(http.StatusBadRequest, "Username should not be empty")
		return
	}
	if len(username) > 255 {
		c.String(http.StatusBadRequest, "Username should not be more than 255 characters")
	}
	password, result := c.GetPostForm("password")
	if !result || password == "" {
		c.String(http.StatusBadRequest, "Password should not be empty")
		return
	}

	err := user.Register(username, password)

	if err != nil {
		c.String(http.StatusOK, "Could not register user. Already registered?")
		log.Error(err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/login")
}

func logoutHandler(c *gin.Context) {
	ctx := c.MustGet("session").(context.Context)

	// Redirect to login if not logged in
	if !sessionManager.Exists(ctx, "user_id") {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	err := sessionManager.Destroy(ctx)
	if err != nil {
		log.Error(err)
		c.String(http.StatusInternalServerError, "An error occured")
		return
	}
	c.Redirect(http.StatusFound, "/")
}