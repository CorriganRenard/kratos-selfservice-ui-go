package handlers

import (
	"log"
	"net/http"

	"github.com/CorriganRenard/kratos-selfservice-ui-go/api_client"
	"github.com/benbjohnson/hashfs"
)

// LoginParams configure the Login http handler
type LoginParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the kratos URL to redirect the browser to,
	// when the user wishes to login, and the 'flow' query param is missing
	FlowRedirectURL string
	CSRFCookieName  string
}

// Login handler displays the login screen
func (lp LoginParams) Login(w http.ResponseWriter, r *http.Request) {

	// Start the login flow with Kratos if required
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		log.Printf("No flow ID found in URL, initializing login flow, redirect to %s", lp.FlowRedirectURL)
		http.Redirect(w, r, lp.FlowRedirectURL, http.StatusMovedPermanently)
		return
	}

	csrfc, err := r.Cookie(lp.CSRFCookieName)
	if err != nil {
		log.Printf("Error getting csrf cookie: %v", err)
		ErrorHandler(w, r, err)
		return
	}

	log.Print("Calling Kratos API to get self service login")
	loginFlow, res, err := api_client.OryClient().FrontendAPI.GetLoginFlow(r.Context()).Cookie(csrfc.String()).Id(flow).Execute()
	if err != nil {
		log.Printf("Error getting self service login flow: %v, redirecting to resetting flow id", err)
		log.Printf("full res: %v", res)
		http.Redirect(w, r, lp.FlowRedirectURL, http.StatusMovedPermanently)
		return
	}

	log.Printf("login flow ui nodes: %#v", loginFlow.Ui.Nodes)

	dataMap := map[string]interface{}{
		"Title":    "Login",
		"Method":   loginFlow.Ui.Method,
		"Action":   loginFlow.Ui.Action,
		"Fields":   loginFlow.Ui.Nodes,
		"Messages": loginFlow.Ui.Messages,
		"flow":     flow,
		//"config":      loginFlow.Ui.GetNodes(),
		"fs":          lp.FS,
		"pageHeading": "Login",
	}
	if err = GetTemplate(loginPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
