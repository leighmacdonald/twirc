package twirc

import "time"

type IRCMessage struct {
	Sent    time.Time
	Message string
}

type Viewer struct {
	Name        string
	HostMask    string
	ConnectedAt time.Time
	History     []string
}

type ActiveViewers struct {
	viewers map[string]*Viewer
}

func (av *ActiveViewers) Add(viewer *Viewer) bool {
	if _, ok := av.viewers[viewer.Name]; !ok {
		//do something here
		av.viewers[viewer.Name] = viewer
		return true
	}
	return false
}

func GetViewer(name string) {

}

func NewViewer(name string, host_mask string) *Viewer {
	return &Viewer{
		Name:        name,
		HostMask:    host_mask,
		ConnectedAt: time.Now(),
	}
}
