package provider

import (
	"time"
)

type Domain struct {
	Name      string    `json:"name"`
	Alias     string    `json:"alias"`
	CacheFile string    `json:"cache_file"`
	CachedAt  time.Time `json:"cached_at"`
	Nodes     []Node    `json:"nodes"`
}

func NewDomain(name string) Domain {
	return Domain{
		Name:      name,
		Alias:     "",
		CacheFile: "",
		Nodes: []Node{
			Node{},
		},
	}
}

func (d *Domain) SetAlias(a string) {
	d.Alias = a
}

func (d *Domain) SetCacheFile(f string) {
	d.CacheFile = f
	d.CachedAt = time.Now()
}

func (d *Domain) IsEqualWithName(other string) bool {
	return d.Name == other
}

func (d *Domain) GetNode(ip string) (Node, bool) {
	for i, n := range d.Nodes {
		if n.IsEqual(ip) {
			return d.Nodes[i], true
		}
	}
	return Node{}, false
}

func (d *Domain) SetNode(node Node) {
	for i, n := range d.Nodes {
		if n.IsEqual(node.IP) {
			d.Nodes[i] = node
		}
	}
	d.Nodes = append(d.Nodes, node)
}

func (d *Domain) GetResolvedIP() (ip []string) {
	for _, n := range d.Nodes {
		ip = append(ip, n.GetIP())
	}
	return
}
