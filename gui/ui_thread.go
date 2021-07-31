package gui

import (
	"github.com/coyim/gotk3adapter/glibi"
)

//TODO: could this use a compiling flag to generate a noop function when released?
func assertInUIThread() {
	if g.glib.MainDepth() == 0 {
		panic("This function has to be called from the UI thread")
	}
}

//GTK process events in glib event loop (see [1]). In order to keep the UI
//responsive, it is a good practice to not block long running tasks in a signal's
//callback (you dont want a button to keep looking pressed for a couple of seconds).
//doInUIThread schedule the function to run in the next
//1 - https://developer.gnome.org/glib/unstable/glib-The-Main-Event-Loop.html
//TODO: Try other patterns and expose them as API. Example: http://www.mono-project.com/docs/gui/gtksharp/responsive-applications/
func doInUIThread(f func()) {
	_ = g.glib.IdleAdd(f)
}

type inUIThread struct {
	g Graphics
}

func (i *inUIThread) assertInUIThread() Graphics {
	return i.g
}

type outsideUIThread struct {
	doInUIThread func(func(*inUIThread))
}

func (*outsideUIThread) assertInUIThread() Graphics {
	panic("This function has to be called from the UI thread")
}

type uiThread interface {
	assertInUIThread() Graphics
}

// FINALIZER FUNCTIONALITY

// finalizerBufferCapability determines how many waiting finalizers we can have before we start blocking
const finalizerBufferCapability = 10000

// finalizerPollRounds determines how many finalizers we will run before giving back control to the UI thread
const finalizerPollRounds = 100

// finalizerPollTime is the time in between rounds of finalizers
const finalizerPollTime = 1000 // milliseconds

func registerFinalizerReaping(gl glibi.Glib) {
	fins := make(chan func(), finalizerBufferCapability)

	gl.SetFinalizerStrategy(func(f func()) {
		fins <- f
	})

	// This will ALWAYS be called in the UI thread
	finalizerPoller := func() bool {
		rounds := 0

		for rounds < finalizerPollRounds {
			select {
			case ff := <-fins:
				ff()
				rounds++
			default:
				return true
			}
		}

		return true
	}

	gl.TimeoutAdd(finalizerPollTime, finalizerPoller)
}
