package handlers

import (
	"log"
	"net/http"

	"github.com/CorriganRenard/kratos-selfservice-ui-go/api_client"
	"github.com/benbjohnson/hashfs"
)

// VerificationParams configures the email verification handler
type VerificationParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the Kratos URL to redirect the browser to when the flow query param is missing
	FlowRedirectURL string
	CSRFCookieName  string
}

// VerifyEmail handler starts the email verification process
func (vp VerificationParams) VerifyEmail(w http.ResponseWriter, r *http.Request) {

	// Start the verification flow with Kratos if required
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		log.Printf("No flow ID found in URL, initializing verification flow, redirect to %s", vp.FlowRedirectURL)
		http.Redirect(w, r, vp.FlowRedirectURL, http.StatusMovedPermanently)
		return
	}

	cookieHeader := r.Header.Get("Cookie")

	log.Print("Calling Kratos API to get self-service verification flow")
	verificationFlow, res, err := api_client.OryClient().FrontendAPI.GetVerificationFlow(r.Context()).Cookie(cookieHeader).Id(flow).Execute()
	if err != nil {
		log.Printf("Error getting self-service verification flow: %v, redirecting to /", err)
		log.Printf("full res: %v", res)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	log.Printf("verification flow ui nodes: %#v", verificationFlow.Ui.Nodes)

	dataMap := map[string]interface{}{
		"Title":       "Email Verification",
		"Method":      verificationFlow.Ui.Method,
		"Action":      verificationFlow.Ui.Action,
		"Fields":      verificationFlow.Ui.Nodes,
		"Messages":    verificationFlow.Ui.Messages,
		"flow":        flow,
		"fs":          vp.FS,
		"pageHeading": "Verify Email",
	}
	if err = GetTemplate(verificationPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
