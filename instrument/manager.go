package instrument

import "turbojet/cli"

const (
	instrumentURL = "https://www.17ce.com/"
)

type Manager struct{}

func Refresh(c *cli.Context, domains []string, url string) ([]map[string]string, error) {
	l := NewLoader()
	w := c.Writer()
	l.Init(c)
	sources, err := l.Process(c, domains, url)
	if err != nil {
		cli.Printf(w, "Loading page source error: %s\n", err)
		return nil, err
	}
	for _, s := range sources {
		cli.Printf(w, "**Refreshed domain: [%s]\n", s["domain"])
	}
	return sources, nil
}
