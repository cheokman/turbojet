package provider

import "time"

const (
	INITED = iota
	HEATING
	HEATED
	INVALID
)

type Node struct {
	IP        string    `json:"IP"`
	Location  string    `json:"location"`
	Alias     string    `json:"alias"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	HeatedAt  time.Time `json:"heated_at"`
}

func NewNode() Node {
	return Node{}
}

func (n *Node) SetIP(value string) {
	n.IP = value
	n.Status = INITED
	n.CreatedAt = time.Now()
}

func (n *Node) SetHeating() {
	n.Status = HEATING
}

func (n *Node) SetHeated() {
	n.Status = HEATED
	n.HeatedAt = time.Now()
}

func (n *Node) IsEqual(ip string) bool {
	if n.IP == ip {
		return true
	}
	return false
}

func (n *Node) GetIP() string {
	return n.IP
}
