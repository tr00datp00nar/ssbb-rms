package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"ssbb-rms/auth"
	"ssbb-rms/database"
	"ssbb-rms/models"

	"github.com/gofiber/fiber/v2/middleware/redirect"

	"github.com/markbates/goth/gothic"
)

// Authentication Handlers

func Auth(res http.ResponseWriter, req *http.Request) {
	// try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		t, _ := template.New("user").Parse(userTemplate)
		t.Execute(res, gothUser)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}

func AuthCallback(res http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	t, _ := template.New("user").Parse(userTemplate)
	t.Execute(res, user)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	gothic.Logout(res, req)
	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>`

func SignUp(res http.ResponseWriter, req *http.Request) {
	// Connect to database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	// Start a new authentication
	auth.NewAuth()
	// Get user authetication data from provider
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Println("Could not authenticate the user.")
		return
	}
	// Check if user already exists in the database
	var dbUsers []models.Users
	db.Table("Users").Scan(&dbUsers)
	for _, dbUser := range dbUsers {
		if user.Email == dbUser.Email {
			redirect.New(redirect.Config{
				Rules: map[string]string{
					"/sign-up": "/login",
				},
				StatusCode: 301,
			})
		}
		// Sign up the user
		// ...
		// ...
		redirect.New(redirect.Config{
			Rules: map[string]string{
				"/sign-up": "/dashboard",
			},
		})
	}
	return
}

func SignIn(res http.ResponseWriter, req *http.Request) {
	// Connect to database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// Start a new authentication
	auth.NewAuth()

	// Try to get the user without reauthenticating
	if _, err := gothic.CompleteUserAuth(res, req); err == nil {
		redirect.New(redirect.Config{
			Rules: map[string]string{
				"/sign-in": "/dashboard",
			},
		})
	} else {
		// Start authentication process
		gothic.BeginAuthHandler(res, req)

		// Get user authetication data from provider
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Println("Could not authenticate the user.")
		}

		// Check if user already exists in the database
		var dbUsers []models.Users
		db.Table("Users").Scan(&dbUsers)
		for _, dbUser := range dbUsers {
			if user.Email != dbUser.Email {
				redirect.New(redirect.Config{
					Rules: map[string]string{
						"/sign-in": "/sign-up",
					},
					StatusCode: 301,
				})
			}
			// Redirect to dashboard
			redirect.New(redirect.Config{
				Rules: map[string]string{
					"/sign-in": "/dashboard",
				},
			})
		}
	}
}
