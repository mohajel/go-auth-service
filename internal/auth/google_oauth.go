package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config

func InitGoogleOAuth() {
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// redirects to Google login page
func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("random-state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handles callback from Google
func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != "random-state" {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	token, err := googleOAuthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	var userInfo map[string]interface{}
	json.Unmarshal(data, &userInfo)

	fmt.Fprintf(w, "✅ Google login success!\n\nUser Info:\n%v", userInfo)
}
