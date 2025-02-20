package handlers

import (
	"log"
	"net/http"

	"github.com/CorriganRenard/kratos-selfservice-ui-go/api_client"
	"github.com/benbjohnson/hashfs"
)

// RecoveryParams configure the Recovery http handler
type RecoveryParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the kratos URL to redirect the browser to,
	// when the user wishes to start recovery, and the 'flow' query param is missing
	FlowRedirectURL string

	CSRFCookieName string
}

// Recovery handler displays the recovery screen, which allows the user to enter
// and email address, the email contains a link to authenticate the user
func (rp RecoveryParams) Recovery(w http.ResponseWriter, r *http.Request) {

	// Start the recovery flow with Kratos if required
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		log.Printf("No flow ID found in URL, initializing recovery flow, redirect to %s", rp.FlowRedirectURL)
		http.Redirect(w, r, rp.FlowRedirectURL, http.StatusMovedPermanently)
		return
	}

	csrfc, err := r.Cookie(rp.CSRFCookieName)
	if err != nil {
		log.Printf("Error getting csrf cookie: %v", err)
		ErrorHandler(w, r, err)
		return
	}

	log.Printf("Calling Kratos API to get self service recovery")
	recoveryFlow, _, err := api_client.OryClient().FrontendAPI.GetRecoveryFlow(r.Context()).Cookie(csrfc.String()).Id(flow).Execute()
	if err != nil {
		log.Printf("Error getting self service recovery flow: %v, redirecting to /", err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	dataMap := map[string]interface{}{
		"Title":    "Recovery",
		"Method":   recoveryFlow.Ui.Method,
		"Action":   recoveryFlow.Ui.Action,
		"Fields":   recoveryFlow.Ui.Nodes,
		"Messages": recoveryFlow.Ui.Messages,
		"State":    recoveryFlow.State,
		"flow":     flow,
		//"config":      recoveryFlow.Ui.GetNodes(),
		"fs":          rp.FS,
		"pageHeading": "Recovery",
	}
	if err = GetTemplate(recoveryPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
	// log.Printf("Recovery state: %v", res.GetPayload().State)
	// dataMap := map[string]interface{}{
	// 	"flow":        flow,
	// 	"link":        res.GetPayload().Methods["link"].Config,
	// 	"state":       res.GetPayload().State,
	// 	"fs":          rp.FS,
	// 	"pageHeading": "Recover your account",
	// }
	// if err = GetTemplate(recoveryPage).Render("layout", w, r, dataMap); err != nil {
	// 	ErrorHandler(w, r, err)
	// }
}
