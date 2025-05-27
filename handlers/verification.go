package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CorriganRenard/kratos-selfservice-ui-go/api_client"
	"github.com/CorriganRenard/kratos-selfservice-ui-go/options"
	"github.com/benbjohnson/hashfs"
)

// VerificationParams configures the email verification handler
type VerificationParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the Kratos URL to redirect the browser to when the flow query param is missing
	FlowRedirectURL string
	CSRFCookieName  string
	Options         *options.Options
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

	// Create a buffer to execute nested templates
	var cardContent template.HTML
	cardTemplate := GetTemplate(verificationPage)
	cardBuffer := &strings.Builder{}
	err = cardTemplate.tmpl.ExecuteTemplate(cardBuffer, "card_content", map[string]interface{}{
		"Method":   verificationFlow.Ui.Method,
		"Action":   verificationFlow.Ui.Action,
		"Fields":   verificationFlow.Ui.Nodes,
		"Messages": verificationFlow.Ui.Messages,
		"flow":     flow,
	})
	if err == nil {
		cardContent = template.HTML(cardBuffer.String())
	}

	var pageContent template.HTML
	pageTemplate := GetTemplate(verificationPage)
	pageBuffer := &strings.Builder{}
	err = pageTemplate.tmpl.ExecuteTemplate(pageBuffer, "page_content", map[string]interface{}{
		"Method":      verificationFlow.Ui.Method,
		"Action":      verificationFlow.Ui.Action,
		"Fields":      verificationFlow.Ui.Nodes,
		"Messages":    verificationFlow.Ui.Messages,
		"flow":        flow,
		"CardContent": cardContent,
	})
	if err == nil {
		pageContent = template.HTML(pageBuffer.String())
	}

	dataMap := map[string]interface{}{
		"Title":       "Email Verification",
		"Method":      verificationFlow.Ui.Method,
		"Action":      verificationFlow.Ui.Action,
		"Fields":      verificationFlow.Ui.Nodes,
		"Messages":    verificationFlow.Ui.Messages,
		"flow":        flow,
		"fs":          vp.FS,
		"pageHeading": "Verify Email",
		"PageContent": pageContent,
		"CardContent": cardContent,
		"siteName":    vp.Options.SiteName,
		"faviconURL":  vp.Options.FaviconURL,
	}
	if err = GetTemplate(verificationPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
