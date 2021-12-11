package thinkutils

import (
	"github.com/howeyc/fsnotify"
)

type OnFileEvent func(szFile string)

type ThinkNotify struct {
	pWatcher *fsnotify.Watcher
	OnCreate OnFileEvent
	OnModify OnFileEvent
	OnDelete OnFileEvent
}

func (this *ThinkNotify) Watch(szPath string) {
	//log.Info("%p", this)
	pWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-pWatcher.Event:
				//log.Info("event:", ev.String())
				if ev.IsCreate() && this.OnCreate != nil {
					if this.OnCreate != nil {
						go this.OnCreate(ev.Name)
					}
				} else if ev.IsAttrib() || ev.IsModify() || ev.IsRename() {
					if this.OnModify != nil {
						go this.OnModify(ev.Name)
					}
				} else if ev.IsDelete() {
					if this.OnDelete != nil {
						go this.OnDelete(ev.Name)
					}
				} else {
					if this.OnModify != nil {
						go this.OnModify(ev.Name)
					}
				}
			case err := <-pWatcher.Error:
				log.Error("error:", err.Error())
			}
		}
	}()

	err = pWatcher.Watch(szPath)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	<-done
	pWatcher.Close()
}
