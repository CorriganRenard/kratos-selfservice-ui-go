package handlers

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/CorriganRenard/kratos-selfservice-ui-go/api_client"
	"github.com/benbjohnson/hashfs"
)

// LogoutParams configure the Logout http handler
type LogoutParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the kratos URL to redirect the browser to,
	// when the user wishes to logout, and the 'flow' query param is missing
	FlowRedirectURL string
	CSRFCookieName  string
}

// Logout handler clears the session & logs the user out
func (lp LogoutParams) Logout(w http.ResponseWriter, r *http.Request) {

	// Start the logout flow with Kratos if required
	// flow := r.URL.Query().Get("flow")
	// if flow == "" {
	// 	log.Printf("No flow ID found in URL, initializing logout flow, redirect to %s", lp.FlowRedirectURL)
	// 	http.Redirect(w, r, lp.FlowRedirectURL, http.StatusMovedPermanently)
	// 	return
	// }
	// forward headers
	cookieHeader := r.Header.Get("Cookie")

	log.Print("Calling Kratos API to get self service logout")
	logoutFlow, res, err := api_client.OryClient().FrontendAPI.CreateBrowserLogoutFlow(r.Context()).Cookie(cookieHeader).Execute()
	if err != nil {
		log.Printf("Error creating self service logout flow: %v, redirecting to /", err)
		log.Printf("full res: %v", res)
		//http.Redirect(w, r, "/", http.StatusMovedPermanently)
		//return
	}

	log.Printf("logout flow url: %v", logoutFlow.LogoutUrl)

	// Redirect user to the post-logout URL
	if logoutFlow.LogoutUrl != "" {
		http.Redirect(w, r, logoutFlow.LogoutUrl, http.StatusFound)
	} else {
		// Fallback redirect in case there's no logout URL
		//http.Redirect(w, r, "/", http.StatusFound)
	}

	dataMap := map[string]interface{}{
		"fs":          lp.FS,
		"pageHeading": "Logged Out",
	}
	if err := GetTemplate(logoutPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
