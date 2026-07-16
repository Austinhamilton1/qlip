package ui

import "time"

type ConnectionMode int

const (
	ClientMode ConnectionMode = iota
	ServerMode
	BothMode
)

type Page int

const (
	PageHome Page = iota
	PageConnection
	PageSettings
)

type SavedConnection struct {
	Name string
	Host string
	Port int
}

type ClientInfo struct {
	Hostname string
	IP       string
	Port     int
}

type State struct {
	Mode            ConnectionMode
	Connected       bool
	ConnectedHost   string
	ConnectedIP     string
	ConnectedPort   int
	LastUpdate      time.Time
	Clients         []ClientInfo
	SavedConnection []SavedConnection
}
