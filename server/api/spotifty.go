package api

import (
	"time"

	"github.com/bmxguy100/spotify"

	log "github.com/sirupsen/logrus"
)

type PlaybackState struct {
	IsAuthenticated bool
	AuthUrl         string
	Err             error
	State           spotify.PlayerState
}

var PlaybackStateChannel = make(chan PlaybackState)

const rateLimit = time.Duration(time.Second * 2)

var cachedStateTime time.Time
var cachedState *spotify.PlayerState

func SpotifyServer() {
	go httpServer()

	log.WithField("URL", baseURL).Info("Please Authenticate")

	for {
		var state spotify.PlayerState
		var err error

		if isAuthenticated {
			elapsedTime := time.Since(cachedStateTime)
			if cachedState == nil || elapsedTime >= rateLimit {
				cachedState, err = client.PlayerState()
				if err == nil && cachedState != nil {
					state = *cachedState
				}
			} else {
				// Use cache and estimate time elapsed
				state = *cachedState
				if state.Playing {
					state.Progress += int(elapsedTime.Milliseconds())
				}
			}
		}

		PlaybackStateChannel <- PlaybackState{
			IsAuthenticated: isAuthenticated,
			AuthUrl:         baseURL,
			State:           state,
			Err:             err,
		}
	}
}
