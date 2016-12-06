package integration_test

import (
	"github.com/draganm/go-reactor"
)

type DataScreen struct {
	Text string
	ctx  reactor.ScreenContext
}

func (d *DataScreen) Mount() {
	d.ctx.UpdateScreen(&reactor.DisplayUpdate{
		Model: &reactor.DisplayModel{
			Element: "div",
			Attributes: map[string]interface{}{
				"className": "top",
			},
			Children: []*reactor.DisplayModel{
				&reactor.DisplayModel{
					Text: d.Text,
				},
			},
		},
	})
}

func (s *DataScreen) Unmount() {
}

func (s *DataScreen) OnUserEvent(*reactor.UserEvent) {

}

func (d *DataScreen) OnText(text string) {
	d.ctx.UpdateScreen(&reactor.DisplayUpdate{
		Model: &reactor.DisplayModel{
			Element: "div",
			Attributes: map[string]interface{}{
				"className": "top",
			},
			Children: []*reactor.DisplayModel{
				&reactor.DisplayModel{
					Text: text,
				},
			},
		},
	})
}

func NewDataScreen(text string, ctx reactor.ScreenContext) *DataScreen {
	return &DataScreen{
		Text: text,
		ctx:  ctx,
	}
}
