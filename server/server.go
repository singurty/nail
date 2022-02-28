package server

import (
	"net/http"
	"context"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	log "github.com/sirupsen/logrus"

	"github.com/singurty/nail/db"
)

var sessionManager *scs.SessionManager

func Start() error {
	// Initialize session manager
	sessionManager = scs.New()
	sessionManager.Cookie.Persist = false;
	sessionManager.Store = pgxstore.New(db.DBpool)
	r := gin.Default()
	r.Use(sessionMiddleware)

	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", rootHandler)
	r.GET("/logout", logoutHandler)

	r.POST("/login", loginHandler)
	r.GET("/login", loginHandler)

	r.POST("/register", registerHandler)
	r.GET("/register", registerHandler)

	r.POST("/post", postHandler)
	r.GET("/create_post", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_post.html", nil)
	})
	
	// Get port from env if present otherwise set to 8080
	port := os.Getenv("PORT")
	_, err := strconv.Atoi(port)
	if port == "" || err != nil {
		port = "8080"
	}

	err = r.Run(":" + port)
	if err != nil {
		return err
	}
	return nil
}

// Set the session context
func sessionMiddleware(c *gin.Context) {
	token, _ := c.Cookie(sessionManager.Cookie.Name)
	ctx, err := sessionManager.Load(c, token)
	if err != nil {
		log.Fatal(err)
	}
	c.Set("session", ctx)
	c.Next()
}

func rootHandler(c *gin.Context) {
	// Redirect to login page if not logged in
	ctx := c.MustGet("session").(context.Context)
	if !sessionManager.Exists(ctx, "user_id") {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "index.html", nil)
}
