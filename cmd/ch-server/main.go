package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"b.yadro.com/sys/ch-server/config"
	"b.yadro.com/sys/ch-server/infrastructure"
	"b.yadro.com/sys/ch-server/infrastructure/repository"
	"b.yadro.com/sys/ch-server/interfaces"
	"b.yadro.com/sys/ch-server/usecases"
)

func main() {
	// Load config from ./config/config.yaml
	config, err := config.NewConfig()
	if err != nil {
		stderr.Fatalf("Unable to open config: %s", err)
	}
	// Create a new tusd composer to bind tusd datastore.
	composer := infrastructure.NewStoreComposer()
	// Create a new tusd handler instance to bind with http server.
	tusdHandler, err := infrastructure.TusdConfig(
		composer,
		config.Tusd.File_path,
		config.Tusd.URL_path)
	if err != nil {
		stderr.Fatalf("Unable to create handler: %s", err)
	}
	// Create a new db handler to save metedata
	dbHandler, err := repository.NewBoltHandler(
		filepath.Join(
			config.DB.File_path,
			config.DB.File_name),
		stderr)
	if err != nil {
		stderr.Fatalf("Unable to create handler: %s", err)
	}
	syrAuth, err := infrastructure.NewSYRAuth(
		config.SYR.URL_auth,
		config.SYR.Token_field,
		config.SYR.Token_header,
		config.SYR_login,
		config.SYR_password)
	if err != nil {
		stderr.Fatalf("Unable to create SYR authorize: %s", err)
	}

	// Create a new SYR handler to upload files
	syrConfig, err := infrastructure.NewSYRConfig(
		config.SYR.URL_path,
		config.SYR.File_path,
		config.SYR.Field_form,
		config.SYR.File_ext)
	if err != nil {
		stderr.Fatalf("Unable to create SYR config: %s", err)
	}
	// Create a new SYR handler to upload files
	syrHandler, err := infrastructure.NewSYRHandler(
		syrConfig,
		syrAuth,
		stderr)
	if err != nil {
		stderr.Fatalf("Unable to create SYR config: %s", err)
	}
	// Create a new handler to invoke tusd functions
	invokeHandler, err := infrastructure.NewTusdInvoke(composer)
	if err != nil {
		stderr.Fatalf("Unable to create handler: %s", err)
	}
	// Create a new agent to manage metadata action and user stores
	repositoryHandler, err := interfaces.NewDbDataRepo(dbHandler, invokeHandler, "root")
	if err != nil {
		stderr.Fatalf("Unable to create handler: %s", err)
	}
	httpClientHandler, err := interfaces.NewHTTPClient(syrHandler, stdout)
	if err != nil {
		stderr.Fatalf("Unable to create handler: %s", err)
	}
	dataAgent, err := usecases.NewDataAgent(
		repositoryHandler,
		httpClientHandler)
	if err != nil {
		stderr.Fatalf("Unable to create dataAgent: %s", err)
	}
	// Create hooksHandler to invoke hooks
	hooksHandler, err := interfaces.NewHooksHandler(dataAgent, syrHandler, stdout)
	if err != nil {
		stderr.Fatalf("Unable to create dataAgent: %s", err)
	}

	// Create a new hooks handler to manage notice from tusd
	hooksTusdHandler, err := infrastructure.NewHooksTusdHandler(
		composer,
		hooksHandler,
		stderr)
	if err != nil {
		stderr.Fatalf("Unable to create hooksTusdHandler: %s", err)
	}

	// tusd service will start listening on and accept request at
	http.Handle(
		config.Tusd.URL_path,
		http.StripPrefix(config.Tusd.URL_path, tusdHandler))
	srv := &http.Server{Addr: config.Tusd.URL_addr, Handler: nil}

	// Create wait group variable for goroutines
	var wg sync.WaitGroup
	// Channel for exit app.
	exit := make(chan bool, 4)
	ctx, cancel := context.WithCancel(context.Background())
	// Tusd server goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Start tusd service to listen
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			stderr.Printf("[tusd] Unable to listen: %s", err)
		}
		exit <- true
	}()
	// Tusd Hooks goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Start hooks from tusd to usecases
		err := hooksTusdHandler.RunHooks(ctx, tusdHandler)
		if err != nil {
			stderr.Printf("[hooks] Unable to run: %s", err)
		}
		exit <- true
	}()
	// SYR client goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Start hooks from tusd to usecases
		err := syrHandler.Run(ctx, nil)
		if err != nil {
			stderr.Printf("[syr] Unable to run: %s", err)
		}
		exit <- true
	}()
	// Handle exit of program
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Setting up signal capturing
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		select {
		case <-stop:
		case <-exit:
		}
		exit <- true
	}()

	<-exit
	ctxtimeout, srvcancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctxtimeout); err != nil {
		stderr.Printf("[tusd] Shutdown: %s", err)
	}
	srvcancel()
	cancel()
	wg.Wait()
}
