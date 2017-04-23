# go-reactor
Framework for writing completely reactive web applications.


## Hello-World
Simple one page hello world program that counts number of clicks:
```go
package main

import (
	"fmt"
	"sync"

	"github.com/draganm/go-reactor"
)

func main() {
	r := reactor.New()
	r.AddScreen("/", indexScreenFactory)
	r.Serve(":8080")
}

func indexScreenFactory(ctx reactor.ScreenContext) reactor.Screen {
	return &indexScreen{
		ctx: ctx,
	}
}

var indexScreenTemplate = reactor.MustParseDisplayModel(`
<div className="well">
	You've clicked <mark id="hw" reportEvents="click">hello world</mark>: <span id="count"/> times
</div>
`)

type indexScreen struct {
	sync.Mutex
	ctx          reactor.ScreenContext
	clickCounter int
}

func (i *indexScreen) Mount() {
	i.clickCounter = 0
	i.render()
}

func (i *indexScreen) render() {
	ui := indexScreenTemplate.DeepCopy()
	ui.SetElementText("count", fmt.Sprintf("%d", i.clickCounter))
	i.ctx.UpdateScreen(&reactor.DisplayUpdate{
		Model: ui,
	})
}

func (i *indexScreen) OnUserEvent(evt *reactor.UserEvent) {
	i.Lock()
	defer i.Unlock()

	if evt.ElementID == "hw" {
		i.clickCounter++
		i.render()
	}
}

func (i *indexScreen) Unmount() {}

```

## Updating dependencies

```sh
npm install
gulp
cd public
go generate
cd ..
```
