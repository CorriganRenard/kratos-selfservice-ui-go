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

// LoginParams configure the Login http handler
type LoginParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the kratos URL to redirect the browser to,
	// when the user wishes to login, and the 'flow' query param is missing
	FlowRedirectURL string
	CSRFCookieName  string
	Options         *options.Options
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

	// Create a buffer to execute nested templates
	var cardContent template.HTML
	cardTemplate := GetTemplate(loginPage)
	cardBuffer := &strings.Builder{}
	err = cardTemplate.tmpl.ExecuteTemplate(cardBuffer, "card_content", map[string]interface{}{
		"Method":   loginFlow.Ui.Method,
		"Action":   loginFlow.Ui.Action,
		"Fields":   loginFlow.Ui.Nodes,
		"Messages": loginFlow.Ui.Messages,
		"flow":     flow,
	})
	if err == nil {
		cardContent = template.HTML(cardBuffer.String())
	}

	var pageContent template.HTML
	pageTemplate := GetTemplate(loginPage)
	pageBuffer := &strings.Builder{}
	err = pageTemplate.tmpl.ExecuteTemplate(pageBuffer, "page_content", map[string]interface{}{
		"Method":      loginFlow.Ui.Method,
		"Action":      loginFlow.Ui.Action,
		"Fields":      loginFlow.Ui.Nodes,
		"Messages":    loginFlow.Ui.Messages,
		"flow":        flow,
		"CardContent": cardContent,
	})
	if err == nil {
		pageContent = template.HTML(pageBuffer.String())
	}

	dataMap := map[string]interface{}{
		"Title":       "Login",
		"Method":      loginFlow.Ui.Method,
		"Action":      loginFlow.Ui.Action,
		"Fields":      loginFlow.Ui.Nodes,
		"Messages":    loginFlow.Ui.Messages,
		"flow":        flow,
		"fs":          lp.FS,
		"pageHeading": "Login",
		"PageContent": pageContent,
		"CardContent": cardContent,
		"siteName":    lp.Options.SiteName,
		"faviconURL":  lp.Options.FaviconURL,
	}
	if err = GetTemplate(loginPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
