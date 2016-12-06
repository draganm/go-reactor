package integration_test

import (
	"github.com/draganm/go-reactor"
)

type ClickEventScreen struct {
	ctx     reactor.ScreenContext
	clicked bool
}

var clickEventScreenUI = reactor.MustParseDisplayModel(`
<div>
  <button reportEvents="click">Click me!</button>
  <div className="status" id="status"></div>
</div>
`)

func (d *ClickEventScreen) Mount() {
	d.render()
}

func (d *ClickEventScreen) render() {
	model := clickEventScreenUI.DeepCopy()
	if d.clicked {
		model.SetElementText("status", "clicked!")
	}

	d.ctx.UpdateScreen(&reactor.DisplayUpdate{Model: model})
}

func (s *ClickEventScreen) Unmount() {
}

func (s *ClickEventScreen) OnUserEvent(evt *reactor.UserEvent) {
	if evt.Type == "click" {
		s.clicked = true
	}
	s.render()
}

func NewClickEventScreen(ctx reactor.ScreenContext) *ClickEventScreen {
	return &ClickEventScreen{ctx: ctx}
}
