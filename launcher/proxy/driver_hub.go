// Package driverhub provides a handler for proxying connections to a Selenium server.
package driverhub

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux/mux"
	"github.com/web_test_launcher/launcher/environments/environment"
	"github.com/web_test_launcher/launcher/errors"
	"github.com/web_test_launcher/launcher/healthreporter"
	"github.com/web_test_launcher/launcher/proxy/webdriver"
	"github.com/web_test_launcher/util/httphelper"
)

const (
	compName   = "WebDriver Hub"
	envTimeout = 5 * time.Minute // some environments such as Android take a long time to start up.
)

type hub struct {
	*mux.Router
	env         environment.Env
	healthyOnce sync.Once
	client      *http.Client

	mu       sync.RWMutex
	sessions map[string]http.Handler
	nextID   int
}

// NewHandler creates a handler for /wd/hub paths that delegates to a WebDriver server instance provided by env.
func NewHandler(env environment.Env) http.Handler {
	h := &hub{
		Router:   mux.NewRouter(),
		env:      env,
		sessions: map[string]http.Handler{},
		client:   &http.Client{},
	}

	h.Path("/wd/hub/session").Methods("POST").HandlerFunc(withContext(h.createSession))
	h.Path("/wd/hub/session").HandlerFunc(unknownMethod)
	h.PathPrefix("/wd/hub/session/{sessionID}").HandlerFunc(h.routeToSession)
	h.PathPrefix("/wd/hub/{command}").HandlerFunc(withContext(h.defaultForward))
	h.PathPrefix("/").HandlerFunc(unknownCommand)

	return h
}

func withContext(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(r.Context(), w, r)
	}
}

func (h *hub) routeToSession(w http.ResponseWriter, r *http.Request) {
	sid := mux.Vars(r)["sessionID"]
	h.mu.RLock()
	session := h.sessions[sid]
	h.mu.RUnlock()

	if session == nil {
		invalidSessionID(w, sid)
		return
	}
	session.ServeHTTP(w, r)
}

func (h *hub) createSession(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log.Print("Creating session\n\n")
	var desired map[string]interface{}

	if err := h.waitForHealthyEnv(ctx); err != nil {
		sessionNotCreated(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sessionNotCreated(w, err)
		return
	}

	j := struct {
		Desired map[string]interface{} `json:"desiredCapabilities"`
	}{}

	if err := json.Unmarshal(body, &j); err != nil {
		sessionNotCreated(w, err)
		return
	}

	if j.Desired == nil {
		sessionNotCreated(w, errors.New(compName, "new session request body missing desired capabilities"))
		return
	}

	h.mu.Lock()
	id := h.nextID
	h.nextID++
	h.mu.Unlock()

	desired, err = h.env.StartSession(ctx, id, j.Desired)
	if err != nil {
		sessionNotCreated(w, err)
		return
	}

	log.Printf("Caps: %+v", desired)
	// TODO(fisherii) parameterize attempts based on browser metadata
	driver, err := webdriver.CreateSession(ctx, "http://"+h.env.WDAddress(ctx)+"/wd/hub/", 3, desired)
	if err != nil {
		if err2 := h.env.StopSession(ctx, id); err2 != nil {
			log.Printf("error stopping session after failing to launch webdriver: %v", err2)
		}
		sessionNotCreated(w, err)
		return
	}

	session, err := createSession(id, h, driver, desired)
	if err != nil {
		sessionNotCreated(w, err)
		return
	}

	h.mu.Lock()
	h.sessions[driver.SessionID()] = session
	h.mu.Unlock()

	respJSON := map[string]interface{}{
		"status":    0,
		"sessionId": driver.SessionID(),
		"value":     driver.Capabilities(),
	}

	bytes, err := json.Marshal(respJSON)
	if err != nil {
		unknownError(w, err)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(bytes)
}

func (h *hub) defaultForward(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.waitForHealthyEnv(ctx); err != nil {
		unknownError(w, err)
		return
	}

	if err := httphelper.Forward(ctx, h.env.WDAddress(ctx), w, r); err != nil {
		unknownError(w, err)
	}
}

func (h *hub) waitForHealthyEnv(ctx context.Context) error {
	h.healthyOnce.Do(func() {
		healthyCtx, cancel := context.WithTimeout(ctx, envTimeout)
		defer cancel()
		// ignore error here as we will call and return Healthy below.
		healthreporter.WaitForHealthy(healthyCtx, h.env)
	})
	return h.env.Healthy(ctx)
}
