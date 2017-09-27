package reactor

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/urfave/negroni"

	"github.com/draganm/go-reactor/path"
	"github.com/draganm/go-reactor/public"
)

type ScreenContext struct {
	Path         string
	ConnectionID string
	Params       map[string]string
	UpdateScreen func(*DisplayUpdate)
}

type Screen interface {
	Mount()
	OnUserEvent(*UserEvent)
	Unmount()
}

type ScreenFactory func(ScreenContext) Screen

type screenMatcher struct {
	matcher path.Matcher
	factory ScreenFactory
}

type Reactor struct {
	sync.RWMutex
	matchers              []screenMatcher
	handlers              []negroni.Handler
	notFoundScreenFactory ScreenFactory
}

func New(handlers ...negroni.Handler) *Reactor {
	return &Reactor{
		notFoundScreenFactory: DefaultNotFoundScreenFactory,
		handlers:              handlers,
	}
}

func (r *Reactor) findScreenFactoryForPath(path string) (ScreenFactory, map[string]string) {
	for _, m := range r.matchers {
		params := m.matcher(path)
		if params != nil {
			return m.factory, params
		}
	}
	return r.notFoundScreenFactory, nil
}

func (re *Reactor) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "GET" && r.URL.Path == "/ws" {

		displayChan := make(chan *DisplayUpdate)
		eventChan := make(chan *UserEvent)

		connectionID := uuid.NewV4().String()

		header := http.Header{}

		c := http.Cookie{Name: "BR", Value: "true"}
		header.Add("Set-Cookie", c.String())

		conn, err := upgrader.Upgrade(rw, r, header)
		if err != nil {
			panic(err)
		}

		go func() {
			currentState := &DisplayUpdate{}
			for displayUpdate := range displayChan {
				if !displayUpdate.DeepEqual(currentState) {
					err := conn.WriteJSON(displayUpdate)
					if err != nil {
						break
					}
					currentState = displayUpdate
				}
			}
		}()

		defer func() {
			close(eventChan)
			close(displayChan)
		}()

		go func() {

			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in screen loop", r)
				}
			}()

			path := "/"

		mainLoop:
			for {
				if path == "" {
					path = "/"
				}
				// TODO lock etc.

				screenFactory, params := re.findScreenFactoryForPath(path)
				updater := NewrateLimitedScreenUpdater(200*time.Millisecond, func(upd *DisplayUpdate) { displayChan <- upd })
				ctx := ScreenContext{
					Path:         path,
					ConnectionID: connectionID,
					Params:       params,
					UpdateScreen: newRemoveDuplicatesScreenUpdater(updater.update),
				}

				currentScreen := screenFactory(ctx)

				if currentScreen != nil {

					currentScreen.Mount()

					for evt := range eventChan {
						if evt.Type == "popstate" {
							newPath := strings.TrimPrefix(evt.Value, "#")
							if newPath != path {
								path = newPath
								currentScreen.Unmount()
								updater.close()
								continue mainLoop
							}
						}
						currentScreen.OnUserEvent(evt)
					}
					currentScreen.Unmount()
					updater.close()
					close(displayChan)
					currentScreen = nil
					return
				}

			}
		}()

		for {
			evt := &UserEvent{}
			err := conn.ReadJSON(evt)
			if err != nil {
				return
			}

			eventChan <- evt

		}

	}
	next.ServeHTTP(rw, r)
}

func (r *Reactor) Serve(bind string) {
	handlers := append(r.handlers, negroni.NewStatic(public.AssetFS()))
	n := negroni.New(append(handlers, r)...)

	n.Run(bind)
}

func (r *Reactor) AddScreen(pathPattern string, factory ScreenFactory) error {
	matcher, err := path.NewMatcher(pathPattern)
	if err != nil {
		return err
	}
	r.Lock()
	defer r.Unlock()
	r.matchers = append(r.matchers, screenMatcher{matcher, factory})
	return nil
}

func (r *Reactor) RemoveScreens() {
	r.Lock()
	defer r.Unlock()
	r.matchers = []screenMatcher{}
}
