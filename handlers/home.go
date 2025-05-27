package handlers

import (
	"net/http"

	"github.com/CorriganRenard/kratos-selfservice-ui-go/options"
	"github.com/CorriganRenard/kratos-selfservice-ui-go/session"
	"github.com/benbjohnson/hashfs"
)

// HomeParams configure the Home http handler
type HomeParams struct {
	session.SessionStore
	FS      *hashfs.FS
	Options *options.Options
}

// Home displays a simple homepage
func (p HomeParams) Home(w http.ResponseWriter, r *http.Request) {
	dataMap := map[string]interface{}{
		"fs":            p.FS,
		"kratosSession": p.GetKratosSession(r),
		"headers":       []string{},
		"pageHeading":   "Homepage",
		"siteName":      p.Options.SiteName,
		"faviconURL":    p.Options.FaviconURL,
	}
	if err := GetTemplate(homePage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
