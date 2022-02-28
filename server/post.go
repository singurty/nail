package server

import (
	"net/http"
	"context"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/singurty/nail/posts"
)

func postHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		ctx := c.MustGet("session").(context.Context)
		title, result := c.GetPostForm("title")
		if !result || title == "" {
			c.String(http.StatusBadRequest, "Title should not be empty")
			return
		}
		body, _ := c.GetPostForm("body")
		// Create post
		var err error
		// Create anonymous post if not logged in
		if !sessionManager.Exists(ctx, "user_id") {
			err = posts.Create(title, body, 0)
		} else {
			user_id := sessionManager.Get(ctx, "user_id").(int)
			err = posts.Create(title, body, user_id)
		}
		if err != nil {
			log.Error(err)
			c.String(http.StatusInternalServerError, "Could not create post.")
		}
	}
}
