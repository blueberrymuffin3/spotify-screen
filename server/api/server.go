package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bmxguy100/spotify"
	"github.com/google/uuid"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

var auth = spotify.NewAuthenticator(callbackURL, spotify.ScopeUserReadPlaybackState)
var state = uuid.New().String()
var hostname, _ = os.Hostname()
var baseURL = fmt.Sprintf("http://%s:%s/", hostname, os.Getenv("PORT"))
var callbackURL = fmt.Sprintf("%scallback/", baseURL)

var client spotify.Client
var isAuthenticated = false

func authCallback(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated {
		token, err := auth.Token(state, r)
		if err != nil {
			http.Error(w, "Couln't get token", http.StatusInternalServerError)
			return
		}

		isAuthenticated = true

		log.Info("Recieved authentication token")
		client = auth.NewClient(token)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func rootCallback(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated {
		http.Redirect(w, r, auth.AuthURL(state), http.StatusFound)
		return
	}

	fmt.Fprint(w, "OK")
}

func httpServer() {
	auth.SetAuthInfo(os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))

	http.HandleFunc("/callback/", authCallback)
	http.HandleFunc("/", rootCallback)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")), nil)
	if err != nil {
		log.Fatal(err)
	}
}
