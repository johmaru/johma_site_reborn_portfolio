package main

import (
	"context"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/Johmaru/johma_site_reborn/src/handlers"
	"github.com/Johmaru/johma_site_reborn/src/middleware"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var DevMode = false

func main() {

	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT environment variable must be set.")
	}

	var conf *firebase.Config

	if os.Getenv("FIRESTORE_EMULATOR_HOST") != "" {
		conf = &firebase.Config{
			ProjectID: projectID,
		}
	} else {
		conf = &firebase.Config{ProjectID: projectID}
	}

	var app *firebase.App
	var err error

	if os.Getenv("FIRESTORE_EMULATOR_HOST") != "" {
		app, err = firebase.NewApp(ctx, conf)
	} else {
		sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
		app, err = firebase.NewApp(ctx, conf, sa)
	}

	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to create auth client: %v", err)
	}

	userHandler := handlers.NewUserHandler(firestoreClient, authClient)

	authMiddleware := middleware.NewAuthMiddleware(authClient)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.GET("/", handleIndex)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	api := router.Group("/api")
	{
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", authMiddleware.AuthRequired(), authMiddleware.OwnerRequired(), userHandler.GetUser)
		api.GET("/check-login", handleCheckIsLogin)
	}

	router.GET("/home", handleHome)
	router.GET("/about", handleAbout)
	router.GET("/nowdev", handleNowDev)
	router.GET("/register", handleRegister)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(ctx *gin.Context) {
	theme, err := ctx.Cookie("theme")
	if err != nil {
		theme = "light"
	}
	lang, err := ctx.Cookie("lang")
	if err != nil {
		lang = "ja"
	}

	_, err = ctx.Cookie("isLogin")
	isLogin := err == nil

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "Johma WebSite",
		"theme":   theme,
		"lang":    lang,
		"devmode": DevMode,
		"isLogin": isLogin,
	})
}

func handleHome(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "home.html", nil)
}

func handleAbout(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "about.html", nil)
}

func handleNowDev(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "nowdev.html", nil)
}

func handleCheckIsLogin(ctx *gin.Context) {
	isLogin, err := ctx.Cookie("isLogin")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"isLogin": false,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"isLogin": isLogin,
	})
}

func handleRegister(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", nil)
}
