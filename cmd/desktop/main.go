// Command desktop is Progen Desktop, a desktop app that generates Spring Boot projects.
package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/sivaprasadreddy/progen/internal/desktop/app"
	"github.com/sivaprasadreddy/progen/internal/desktop/config"
	"github.com/sivaprasadreddy/progen/internal/desktop/generator"
	"github.com/sivaprasadreddy/progen/internal/desktop/store"
)

func main() {
	// Templates and defaults live in the user data directory, e.g.
	// ~/Library/Application Support/Progen Desktop/templates.json.
	userData := func() (string, error) { return mygo.App.Path(mygo.PathUserData) }

	mygo.App.WhenReady(func() {
		var win *mygo.Window
		// Generation runs Maven or Gradle, so it runs off the UI thread and
		// reports back on it.
		async := func(work func() func()) {
			go func() {
				apply := work()
				win.Update(apply)
			}()
		}
		service := &generator.Service{}
		prepare := func(c config.GeneratorConfiguration) (app.PreparedProject, error) {
			prepared, err := service.Prepare(c)
			if err != nil {
				return nil, err
			}
			return prepared, nil
		}
		wizard := app.NewWizard(store.NewTemplateService(userData), store.NewDefaultsService(userData), service.GenerateProject, prepare, async)
		// A previewed project waits in a temporary directory until saved.
		mygo.App.OnQuit(wizard.Close)
		win = mygo.NewWindow(mygo.WindowOptions{
			Title:     "Progen Desktop",
			Width:     1200,
			Height:    800,
			MinWidth:  900,
			MinHeight: 640,
			// The sidebar runs under the window buttons and drags the window.
			TitleBarStyle: mygo.TitleBarHiddenInset,
			// Opens where the user left it last time.
			StateKey: "main",
			// The window shows the interface MyGo draws, not a web page.
			Content: ui.View(wizard.View),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
