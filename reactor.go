package reactor

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
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

func (r *Reactor) Serve(bind string) {
	handlers := append(r.handlers, negroni.NewStatic(public.AssetFS()))
	n := negroni.New(handlers...)

	router := httprouter.New()
	router.HandlerFunc("GET", "/ws", newReactorHandler(func(uc chan *DisplayUpdate, ue chan *UserEvent, req *http.Request, id string) http.Header {

		go func() {

			path := "/"

		mainLoop:
			for {
				if path == "" {
					path = "/"
				}
				// TODO lock etc.

				screenFactory, params := r.findScreenFactoryForPath(path)
				updater := NewrateLimitedScreenUpdater(200*time.Millisecond, func(upd *DisplayUpdate) { uc <- upd })
				ctx := ScreenContext{
					Path:         path,
					ConnectionID: id,
					Params:       params,
					UpdateScreen: newRemoveDuplicatesScreenUpdater(updater.update),
				}

				currentScreen := screenFactory(ctx)

				if currentScreen != nil {

					currentScreen.Mount()

					for evt := range ue {
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
					return
				}

			}
		}()

		header := http.Header{}

		c := http.Cookie{Name: "BR", Value: "true"}
		header.Add("Set-Cookie", c.String())

		return header
	}))
	n.UseHandler(router)
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
