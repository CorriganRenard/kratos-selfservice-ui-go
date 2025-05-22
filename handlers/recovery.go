package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strings"

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

	// Create a buffer to execute nested templates
	var cardContent template.HTML
	cardTemplate := GetTemplate(recoveryPage)
	cardBuffer := &strings.Builder{}
	err = cardTemplate.tmpl.ExecuteTemplate(cardBuffer, "card_content", map[string]interface{}{
		"Method":   recoveryFlow.Ui.Method,
		"Action":   recoveryFlow.Ui.Action,
		"Fields":   recoveryFlow.Ui.Nodes,
		"Messages": recoveryFlow.Ui.Messages,
		"State":    recoveryFlow.State,
		"flow":     flow,
	})
	if err == nil {
		cardContent = template.HTML(cardBuffer.String())
	}

	var pageContent template.HTML
	pageTemplate := GetTemplate(recoveryPage)
	pageBuffer := &strings.Builder{}
	err = pageTemplate.tmpl.ExecuteTemplate(pageBuffer, "page_content", map[string]interface{}{
		"Method":      recoveryFlow.Ui.Method,
		"Action":      recoveryFlow.Ui.Action,
		"Fields":      recoveryFlow.Ui.Nodes,
		"Messages":    recoveryFlow.Ui.Messages,
		"State":       recoveryFlow.State,
		"flow":        flow,
		"CardContent": cardContent,
	})
	if err == nil {
		pageContent = template.HTML(pageBuffer.String())
	}

	dataMap := map[string]interface{}{
		"Title":       "Recovery",
		"Method":      recoveryFlow.Ui.Method,
		"Action":      recoveryFlow.Ui.Action,
		"Fields":      recoveryFlow.Ui.Nodes,
		"Messages":    recoveryFlow.Ui.Messages,
		"State":       recoveryFlow.State,
		"flow":        flow,
		"fs":          rp.FS,
		"pageHeading": "Recovery",
		"PageContent": pageContent,
		"CardContent": cardContent,
	}
	if err = GetTemplate(recoveryPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
