package integration_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/draganm/go-reactor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var agoutiDriver *agouti.WebDriver

var r *reactor.Reactor

var _ = BeforeSuite(func(done Done) {
	r = reactor.New()
	go r.Serve(":14344")
	fmt.Println("started server")
	for {
		_, err := http.Get("http://localhost:14344/")
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	agoutiDriver = agouti.ChromeDriver()

	Expect(agoutiDriver.Start()).To(Succeed())
	close(done)
}, 4.0)

var _ = AfterSuite(func(done Done) {
	Expect(agoutiDriver.Stop()).To(Succeed())
	close(done)
})

var _ = Describe("GoReactor", func() {
	var page *agouti.Page

	BeforeEach(func() {
		var err error
		page, err = agoutiDriver.NewPage(agouti.Desired(agouti.Capabilities{
			"chromeOptions": map[string][]string{
				"args": []string{
					"headless",
					// There is no GPU on our Ubuntu box!
					"disable-gpu",

					// Sandbox requires namespace permissions that we don't have on a container
					"no-sandbox",
				},
			},
		}))
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		r.RemoveScreens()
	})

	AfterEach(func(done Done) {
		Expect(page.Destroy()).To(Succeed())
		close(done)
	})

	Context("When there are no screen matchers", func() {
		Context("When I visit the index page", func() {
			It("Shows 'Connected!' message", func(done Done) {
				Expect(page.Navigate("http://localhost:14344/")).To(Succeed())
				Expect(page.First(".page-header")).To(HaveText("Not Found something went wrong"))
				close(done)
			})
		})
	})

	Context("When there is a simple index screen", func() {
		BeforeEach(func() {
			r.AddScreen("/", func(ctx reactor.ScreenContext) reactor.Screen {
				return NewSimpleScreen("This is a test!", ctx)
			})
		})

		It("Shows 'Hello!' message", func(done Done) {
			Expect(page.Navigate("http://localhost:14344/")).To(Succeed())
			Expect(page.First(".top")).To(HaveText("This is a test!"))
			close(done)
		})

	})

	Context("When there is a screen that changes display based on input channel", func() {

		var screen *DataScreen

		BeforeEach(func() {
			r.AddScreen("/", func(ctx reactor.ScreenContext) reactor.Screen {
				screen = NewDataScreen("This is a test!", ctx)
				return screen
			})
		})

		Context("When I visit the page", func() {
			BeforeEach(func() {
				Expect(page.Navigate("http://localhost:14344/")).To(Succeed())
			})

			It("Should show the original message message", func(done Done) {
				Eventually(page.First(".top")).Should(HaveText("This is a test!"))
				close(done)
			})

			Context("When I change the text", func() {
				BeforeEach(func(done Done) {
					screen.OnText("Something else")
					close(done)
				})

				It("Should show the new message", func(done Done) {
					Eventually(page.First(".top")).Should(HaveText("Something else"))
					close(done)
				})

			})

		})

	})

	Context("When there is a screen with a button that reacts to user clicks", func() {

		BeforeEach(func() {
			r.AddScreen("/", func(ctx reactor.ScreenContext) reactor.Screen {
				return NewClickEventScreen(ctx)
			})
			Expect(page.Navigate("http://localhost:14344/")).To(Succeed())
		})

		It("Should have empty status text", func() {
			Expect(page.First(".status")).To(HaveText(""))
		})

		Context("When I click on the button", func() {
			BeforeEach(func(done Done) {
				err := page.FirstByButton("Click me!").Click()
				Expect(err).ToNot(HaveOccurred())
				close(done)
			})

			It("Should have new status text", func() {
				Eventually(page.First(".status")).Should(HaveText("clicked!"))
			})

		})

	})

	Describe("Switching screens", func() {
		Context("When there are two screens registered", func() {
			BeforeEach(func() {
				r.AddScreen("/", func(ctx reactor.ScreenContext) reactor.Screen {
					return NewSimpleScreen("Root Screen", ctx)
				})
				r.AddScreen("/s1", func(ctx reactor.ScreenContext) reactor.Screen {
					return NewSimpleScreen("Subscreen", ctx)
				})
			})
			Context("When I visit the first screen", func() {
				BeforeEach(func() {
					Expect(page.Navigate("http://localhost:14344/")).To(Succeed())
				})
				It("Should show the first screen", func() {
					Expect(page.First(".top")).To(HaveText("Root Screen"))
				})
				Context("When I visit the second screen", func() {
					BeforeEach(func() {
						Expect(page.FirstByLink("clickMe").Click()).To(Succeed())
					})
					It("Should show the first screen", func() {
						Expect(page.First(".top")).To(HaveText("Subscreen"))
					})
				})
			})
			Context("When I visit the second screen", func() {
				BeforeEach(func() {
					Expect(page.Navigate("http://localhost:14344/#/s1")).To(Succeed())
				})
				It("Should show the first screen", func() {
					Expect(page.First(".top")).To(HaveText("Subscreen"))
				})
			})
		})
	})

})
