package internal

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

var buildNumber = "dev"

// Run starts our server
func Run() {
	var port = "8080"
	if v, ok := os.LookupEnv("PORT"); ok {
		port = v
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      newServer(),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorln(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c // wait for this signal

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	defer log.Println("Bye!")
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorln(err)
	}
}
