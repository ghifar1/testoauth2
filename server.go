package main

import (
	"fmt"
	"net/http"
	"testoauth2/oauth_manager"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

func main() {
	_, clientStore, srv := oauth_manager.OauthInit()

	r := gin.Default()

	r.GET("/credentials", func(ctx *gin.Context) {
		CredsRoute(ctx, clientStore)
	})

	g := r.Group("/auth")

	g.Use(func(ctx *gin.Context) {
		Middleware(ctx, srv)
	})

	g.GET("/protected", func(ctx *gin.Context) {
		ProtectedRoute(ctx)
	})

	r.GET("/token", func(ctx *gin.Context) {
		TokenRoute(ctx, srv)
	})

	r.Run()
}

func Middleware(c *gin.Context, srv *server.Server) {
	_, err := srv.ValidationBearerToken(c.Request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid token",
		})
	}

	c.Next()
}

func CredsRoute(c *gin.Context, cs *store.ClientStore) {
	clientId := uuid.New().String()[:8]
	clientSecret := uuid.New().String()[:8]

	err := cs.Set(clientId, &models.Client{
		ID:     clientId,
		Secret: clientSecret,
		Domain: "http://localhost:9094",
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"client_id":     clientId,
		"client_secret": clientSecret,
	})
}

func TokenRoute(c *gin.Context, srv *server.Server) {
	srv.HandleTokenRequest(c.Writer, c.Request)
}

func ProtectedRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
