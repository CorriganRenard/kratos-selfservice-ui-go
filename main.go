package main

import (
	"context"
	"crypto/sha256"
	"embed"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/CorriganRenard/kratos-selfservice-ui-go/api_client"
	"github.com/CorriganRenard/kratos-selfservice-ui-go/handlers"
	"github.com/CorriganRenard/kratos-selfservice-ui-go/middleware"
	"github.com/CorriganRenard/kratos-selfservice-ui-go/options"
	"github.com/CorriganRenard/kratos-selfservice-ui-go/session"

	"github.com/benbjohnson/hashfs"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// staticFS holds the static files, CSS images etc.
// Its baked into the application executable using the embed API - see https://golang.org/pkg/embed/
//
//go:embed static
var staticFS embed.FS

func main() {

	opt := options.NewOptions().SetFromCommandLine()
	if err := opt.Validate(); err != nil {
		log.Fatalf("Error parsing command line: %v", err)
	}
	log.Printf("KratosAdminURL: %s", opt.KratosAdminURL.String())
	log.Printf("KratosPublicURL: %s", opt.KratosPublicURL.String())
	h := sha256.New()
	h.Write([]byte(opt.KratosPublicURL.String()))
	log.Printf("public hash (and cookie name): %x", h.Sum(nil))
	log.Printf("KratosBrowserURL: %s", opt.KratosPublicURL.String())
	log.Printf("BaseURL: %s", opt.BaseURL.String())
	log.Printf("Address: %s", opt.Address())
	log.Printf("Number of Cookie store keys: %d", len(opt.CookieStoreKeyPairs))

	// Setup Kratos API client
	api_client.InitOryClient(*opt.KratosPublicURL)
	api_client.InitOryAdminClient(*opt.KratosAdminURL)

	// Setup sesssion store in cookies
	var store = sessions.NewCookieStore(opt.CookieStoreKeyPairs...)

	// Static assets are wrapped in a hash fs that allows for aggressive http caching
	var fsys = hashfs.NewFS(staticFS)

	// Auth handled by Oathkeeper
	r := mux.NewRouter()

	r.Use(gh.RecoveryHandler(gh.PrintRecoveryStack(true)),
		middleware.NoCacheMiddleware)

	homeP := handlers.HomeParams{
		SessionStore: session.SessionStore{Store: store},
		FS:           fsys,
		Options:      opt,
	}
	r.HandleFunc("/", homeP.Home)

	regP := handlers.RegistrationParams{
		FlowRedirectURL: opt.RegistrationURL(),
		FS:              fsys,
		Options:         opt,
	}
	r.HandleFunc("/registration", regP.Registration)

	settingsP := handlers.SettingsParams{
		FlowRedirectURL: opt.SettingsURL(),
		FS:              fsys,
		Options:         opt,
	}
	r.HandleFunc("/settings", settingsP.Settings)

	loginP := handlers.LoginParams{
		FlowRedirectURL: opt.LoginFlowURL(),
		CSRFCookieName:  "",
		FS:              fsys,
		Options:         opt,
	}
	r.HandleFunc("/login", loginP.Login).Name("login")

	verificationP := handlers.VerificationParams{
		FlowRedirectURL: opt.VerificationFlowURL(),
		CSRFCookieName:  "",
		FS:              fsys,
		Options:         opt,
	}
	r.HandleFunc("/verification", verificationP.VerifyEmail).Name("verification")

	logoutP := handlers.LogoutParams{
		FlowRedirectURL: opt.LogoutFlowURL(),
		CSRFCookieName:  "",
		FS:              fsys,
		Options:         opt,
	}
	r.HandleFunc("/logout", logoutP.Logout)

	recoverP := handlers.RecoveryParams{
		FlowRedirectURL: opt.RecoveryFlowURL(),
		FS:              fsys,
		Options:         opt,
	}
	r.HandleFunc("/recovery", recoverP.Recovery)

	r.PathPrefix("/static/").Handler(hashfs.FileServer(fsys))

	// Following routes must be authenticated, so they get extra middleware
	// using oathkeeper instead
	// authP := middleware.KratosAuthParams{
	// 	SessionStore:      session.SessionStore{Store: store},
	// 	WhoAmIURL:         opt.WhoAmIURL(),
	// 	RedirectUnauthURL: MustURL(r.Get("login")).String(),
	// }

	dashP := handlers.DashboardParams{
		SessionStore: session.SessionStore{Store: store},
		FS:           fsys,
		Options:      opt,
	}
	r.Handle("/dashboard",
		http.HandlerFunc(dashP.Dashboard),
	)

	// Wrap everything in a logger
	logR := gh.LoggingHandler(os.Stdout, r)

	// Start server
	srv := &http.Server{
		Addr: opt.Address(),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      logR, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), opt.ShutdownWait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

// MustURL returns a 'named' URL or panics
func MustURL(r *mux.Route, pairs ...string) *url.URL {
	url, err := r.URL(pairs...)
	if err != nil {
		log.Fatalf("Error r.URL failed with error: %v", err)
	}
	return url
}

// Middleware (this function) makes adding more than one layer of middleware easy
// by specifying them as a list. It will run the last specified handler first.
func Middleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}
