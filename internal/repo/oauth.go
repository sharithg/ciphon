package repo

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v62/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/palantir/go-githubapp/oauth2"
)

func RegisterOAuth2Handler(c githubapp.Config) {
	http.Handle("/api/auth/github", oauth2.NewHandler(
		oauth2.GetConfig(c, []string{"user:email"}),
		// force generated URLs to use HTTPS; useful if the app is behind a reverse proxy
		oauth2.ForceTLS(true),
		// set the callback for successful logins
		oauth2.OnLogin(func(w http.ResponseWriter, r *http.Request, login *oauth2.Login) {
			// look up the current user with the authenticated client
			client := github.NewClient(login.Client)
			user, _, err := client.Users.Get(r.Context(), "")
			// handle error, save the user, ...

			fmt.Println(user, err)

			// redirect the user back to another page
			http.Redirect(w, r, "/dashboard", http.StatusFound)
		}),
	))
}
