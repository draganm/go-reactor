package integration_test

import (
	"github.com/draganm/go-reactor"
)

type SimpleScreen struct {
	text string
	ctx  reactor.ScreenContext
}

func (s SimpleScreen) Mount() {
	mod := simpleScreenUI.DeepCopy()
	mod.SetElementText("text", s.text)
	s.ctx.UpdateScreen(&reactor.DisplayUpdate{Model: mod})
}

var simpleScreenUI = reactor.MustParseDisplayModel(`
	<div>
		<div className="top" id="text"></div>
		<a href="#/s1">clickMe</a>
	</div>
`)

func (s SimpleScreen) Unmount() {
}

func (s SimpleScreen) OnUserEvent(*reactor.UserEvent) {

}

func NewSimpleScreen(text string, ctx reactor.ScreenContext) *SimpleScreen {
	return &SimpleScreen{
		text: text,
		ctx:  ctx,
	}
}
