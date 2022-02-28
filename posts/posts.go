package posts

import (
	"context"
	
	"github.com/jackc/pgconn"
	"github.com/singurty/nail/db"
	log "github.com/sirupsen/logrus"
)

func Create(title, body string, author int) error {
	var err error
	// Anonymous post
	if author == 0 {
		_, err = db.DBpool.Exec(context.Background(), "INSERT INTO posts(title, body) VALUES ($1, $2);", title, body)
		if err != nil {
			return err
		}
	} else {
		_, err = db.DBpool.Exec(context.Background(), "INSERT INTO posts(author, title, body) VALUES ($1, $2, $3);", author, title, body)
		if err != nil {
			return err
		}
	}
	return nil
}
