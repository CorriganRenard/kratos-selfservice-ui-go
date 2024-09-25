package handlers

import (
	_ "embed"
	"io"
	"log"
	"net/http"

	"github.com/CorriganRenard/kratos-selfservice-ui-go/api_client"
	"github.com/benbjohnson/hashfs"
)

// RegistrationParams configure the Registration http handler
type RegistrationParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the kratos URL to redirect the browser to,
	// when the user wishes to register, and the 'flow' query param is missing
	FlowRedirectURL string
	CSRFCookieName  string
}

// Registration directs the user to a page where they can sign-up or
// register to use the site
func (rp RegistrationParams) Registration(w http.ResponseWriter, r *http.Request) {

	// Start the login flow with Kratos if required
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		log.Printf("No flow ID found in URL, initializing registration flow, redirect to %s", rp.FlowRedirectURL)
		http.Redirect(w, r, rp.FlowRedirectURL, http.StatusMovedPermanently)
		return
	}

	csrfc, err := r.Cookie(rp.CSRFCookieName)
	if err != nil {
		log.Printf("Error getting csrf cookie: %v", err)
		ErrorHandler(w, r, err)
		return
	}

	log.Printf("Calling Kratos API to get self service registration")
	registrationFlow, rbody, err := api_client.OryClient().FrontendAPI.GetRegistrationFlow(r.Context()).Cookie(csrfc.String()).Id(flow).Execute()
	if err != nil {
		log.Printf("Error getting self service registration flow: %v, redirecting to /", err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	b, err := io.ReadAll(rbody.Body)
	if err != nil {
		log.Printf("Error getting response body", err)

	}
	log.Printf("flow body: %s", string(b))

	dataMap := map[string]interface{}{
		"Title":    "Registration",
		"Method":   registrationFlow.Ui.Method,
		"Action":   registrationFlow.Ui.Action,
		"Fields":   registrationFlow.Ui.Nodes,
		"Messages": registrationFlow.Ui.Messages,
		"State":    registrationFlow.State,
		"flow":     flow,
		//"config":      registrationFlow.Ui.GetNodes(),
		"fs":          rp.FS,
		"pageHeading": "Registration",
	}
	if err = GetTemplate(registrationPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}

	// // Call Kratos to retrieve the login form
	// params := public.NewGetSelfServiceRegistrationFlowParams()
	// params.SetID(flow)
	// log.Print("Calling Kratos API to get self service registration")
	// res, err := api_client.PublicClient().Public.GetSelfServiceRegistrationFlow(params)
	// if err != nil {
	// 	log.Printf("Error getting self service registration flow %v, redirecting to /", err)
	// 	http.Redirect(w, r, "/", http.StatusMovedPermanently)
	// 	return
	// }
	// dataMap := map[string]interface{}{
	// 	"config":      res.GetPayload().Methods["password"].Config,
	// 	"flow":        flow,
	// 	"fs":          rp.FS,
	// 	"pageHeading": "Registration",
	// }

	// if err = GetTemplate(registrationPage).Render("layout", w, r, dataMap); err != nil {
	// 	ErrorHandler(w, r, err)
	// }
}
