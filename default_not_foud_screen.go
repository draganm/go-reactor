package reactor

type DefaultNotFoundScreen struct {
	ctx ScreenContext
}

var defaultNotFoundScreenUI = MustParseDisplayModel(`
  <bs.PageHeader>Not Found <small>something went wrong</small></bs.PageHeader>
`)

func (d *DefaultNotFoundScreen) Mount() {
	d.ctx.UpdateScreen(&DisplayUpdate{Model: defaultNotFoundScreenUI})
}

func (d *DefaultNotFoundScreen) OnUserEvent(*UserEvent) {

}

func (d *DefaultNotFoundScreen) Unmount() {
}

func DefaultNotFoundScreenFactory(ctx ScreenContext) Screen {
	return &DefaultNotFoundScreen{ctx}
}
