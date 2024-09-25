package handlers

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/CorriganRenard/kratos-selfservice-ui-go/api_client"
	"github.com/benbjohnson/hashfs"
)

// SettingsParams configure the Settings http handler
type SettingsParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the kratos URL to redirect the browser to,
	// when the user wishes to edit their settings, and the 'flow' query param is missing
	FlowRedirectURL string

	CSRFCookieName string
}

// Settings handler displays the Settings screen that allows the user to change their details
func (lp SettingsParams) Settings(w http.ResponseWriter, r *http.Request) {

	// Start the Settings flow with Kratos if required
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		log.Printf("No flow ID found in URL, initializing Settings flow, redirect to %s", lp.FlowRedirectURL)
		http.Redirect(w, r, lp.FlowRedirectURL, http.StatusMovedPermanently)
		return
	}

	csrfc := r.Header.Get("Cookie")

	log.Print("Calling Kratos API to get self service settings")
	settingsFlow, _, err := api_client.OryClient().FrontendAPI.GetSettingsFlow(r.Context()).Cookie(csrfc).Id(flow).Execute()
	if err != nil {
		log.Printf("Error getting self service settings flow: %v, redirecting to /", err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	dataMap := map[string]interface{}{
		"Title":    "Update Profile",
		"Method":   settingsFlow.Ui.Method,
		"Action":   settingsFlow.Ui.Action,
		"Fields":   settingsFlow.Ui.Nodes,
		"Messages": settingsFlow.Ui.Messages,
		"State":    settingsFlow.State,
		"flow":     flow,
		//"config":      settingsFlow.Ui.GetNodes(),
		"fs":          lp.FS,
		"pageHeading": "settings",
	}
	if err = GetTemplate(settingsPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
