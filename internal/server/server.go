package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quan-to/chevron/internal/agent"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/server/pages"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/internal/vaultManager"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
	"io/ioutil"
	"net/http"
)

// GenRemoteSignerServerMux generates a remote signer HTTP Router
func GenRemoteSignerServerMux(slog slog.Instance, sm interfaces.SecretsManager, gpg interfaces.PGPManager) *mux.Router {
	var vm *vaultManager.VaultManager
	log := slog.Scope("MUX")

	if config.VaultStorage {
		vm = vaultManager.MakeVaultManager(log, config.KeyPrefix)
	}

	ge := MakeGPGEndpoint(log, sm, gpg)
	ie := MakeInternalEndpoint(log, sm, gpg)
	te := MakeTestsEndpoint(log, vm)
	kre := MakeKeyRingEndpoint(log, sm, gpg)
	sks := MakeSKSEndpoint(log, sm, gpg)
	tm := agent.MakeTokenManager(log)
	am := agent.MakeAuthManager(log)
	ap := MakeAgentProxy(log, gpg, tm)
	sGql := MakeStaticGraphiQL(log)
	agentAdmin := MakeAgentAdmin(log, tm, am)
	jfc := MakeJFCEndpoint(log, sm, gpg)

	if ge == nil || ie == nil || te == nil || kre == nil || sks == nil || tm == nil || am == nil || ap == nil || agentAdmin == nil {
		slog.Error("One or more services has not been initialized.")
		slog.Error("    GPG Endpoint: %p", ge)
		slog.Error("    Internal Endpoint: %p", ie)
		slog.Error("    Tests Endpoint: %p", te)
		slog.Error("    KeyRing Endpoint: %p", kre)
		slog.Error("    SKS Endpoint: %p", sks)
		slog.Error("    Token Manager: %p", tm)
		slog.Error("    Auth Manager: %p", am)
		slog.Error("    Agent Proxy: %p", ap)
		slog.Error("    Agent Admin: %p", agentAdmin)
		slog.Fatal("Please check if the settings are correct.")
	}

	r := mux.NewRouter()
	// Add for /
	AddHKPEndpoints(log, r.PathPrefix("/pks").Subrouter())
	ge.AttachHandlers(r.PathPrefix("/gpg").Subrouter())
	ie.AttachHandlers(r.PathPrefix("/__internal").Subrouter())
	te.AttachHandlers(r.PathPrefix("/tests").Subrouter())
	kre.AttachHandlers(r.PathPrefix("/keyRing").Subrouter())
	sks.AttachHandlers(r.PathPrefix("/sks").Subrouter())
	jfc.AttachHandlers(r.PathPrefix("/fieldCipher").Subrouter())

	// Add for /remoteSigner
	AddHKPEndpoints(log, r.PathPrefix("/remoteSigner/pks").Subrouter())
	ge.AttachHandlers(r.PathPrefix("/remoteSigner/gpg").Subrouter())
	ie.AttachHandlers(r.PathPrefix("/remoteSigner/__internal").Subrouter())
	te.AttachHandlers(r.PathPrefix("/remoteSigner/tests").Subrouter())
	kre.AttachHandlers(r.PathPrefix("/remoteSigner/keyRing").Subrouter())
	sks.AttachHandlers(r.PathPrefix("/remoteSigner/sks").Subrouter())
	jfc.AttachHandlers(r.PathPrefix("/remoteSigner/fieldCipher").Subrouter())

	// Agent
	ap.AddHandlers(r.PathPrefix("/agent").Subrouter())

	// Agent Admin
	agentAdmin.AddHandlers(r.PathPrefix("/agentAdmin").Subrouter())

	// Static GraphiQL
	sGql.AttachHandlers(r.PathPrefix("/graphiql").Subrouter())

	pages.AddHandlers(r.PathPrefix("/assets").Subrouter())

	// Catch All for unhandled endpoints
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		InitHTTPTimer(log, r)
		CatchAllRouter(w, r, log)
	})

	return r
}

// RunRemoteSignerServer runs a remote signer server asynchronously and returns a stop channel
func RunRemoteSignerServer(slog slog.Instance, sm interfaces.SecretsManager, gpg interfaces.PGPManager) chan bool {

	r := GenRemoteSignerServerMux(slog, sm, gpg)

	listenAddr := fmt.Sprintf("0.0.0.0:%d", config.HttpPort)

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r, // Pass our instance of gorilla/mux in.
	}

	stopChannel := make(chan bool)

	go func() {
		<-stopChannel
		slog.Info("Received STOP. Closing server")
		_ = srv.Close()
		stopChannel <- true
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error(err)
		}
		slog.Info("HTTP Server Closed")
	}()

	slog.Info("Remote Signer is now listening at %s", listenAddr)

	return stopChannel
}

// RunRemoteSignerServerSingleKey runs a single key instance of remote signer server asynchronously and returns a stop channel
func RunRemoteSignerServerSingleKey(slog slog.Instance, sm interfaces.SecretsManager, gpg interfaces.PGPManager) (chan bool, error) {
	slog.Info("Running in single-key mode")

	slog.Info("Loading key from %q", config.SingleKeyPath)
	keyData, err := ioutil.ReadFile(config.SingleKeyPath)

	if err != nil {
		return nil, fmt.Errorf("error reading key file at %s: %q", config.SingleKeyPath, err)
	}

	ctx := context.Background()

	n, err := gpg.LoadKey(ctx, string(keyData))

	if err != nil {
		return nil, fmt.Errorf("error opening private key: %q", err)
	}

	if n == 0 {
		return nil, fmt.Errorf("key parsed sucessfully but no private keys found. check if SINGLE_KEY_PATH points to a private key")
	}

	fps, _ := tools.GetFingerPrintsFromKey(string(keyData))

	fp := fps[0]

	slog.Info("Key loaded. Unlocking key %q", fp)

	err = gpg.UnlockKey(ctx, fp, config.SingleKeyPassword)

	if err != nil {
		return nil, fmt.Errorf("cannot unlock key %q with password provided by SINGLE_KEY_PASSWORD environment", fp)
	}

	slog.Info("Key unlocked. Setting default Agent Key Fingerprint to %q", fp)
	config.AgentKeyFingerPrint = fp

	r := GenRemoteSignerServerMux(slog, sm, gpg)

	listenAddr := fmt.Sprintf("0.0.0.0:%d", config.HttpPort)

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r, // Pass our instance of gorilla/mux in.
	}

	stopChannel := make(chan bool)

	go func() {
		<-stopChannel
		slog.Info("Received STOP. Closing server")
		_ = srv.Close()
		stopChannel <- true
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error(err)
		}
		slog.Info("HTTP Server Closed")
	}()

	slog.Info("Remote Signer is now listening at %s", listenAddr)

	return stopChannel, nil
}
