package instrument

import (
	"encoding/json"
	"fmt"
	"time"
	"turbojet/cli"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	seleniumPath        = "driver/selenium-server-standalone-3.141.59.jar"
	defaultLoaderDriver = "chrome"
	chromeDriverPath    = "driver/chromedriver"
	port                = 8080
	timeToWaitLoading   = 40
)

type Loader struct {
	HeadLess bool
	Driver   string
	caps     selenium.Capabilities
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Init(c *cli.Context) {
	l.caps = selenium.Capabilities{
		"browserName": "chrome",
	}

	imgCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imgCaps,
		Path:  "",
		Args: []string{
			"--headless",
			"--start-maximized",
			//"--window-size=1200x600",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
			"--disable-gpu",
			"--disable-impl-side-painting",
			"--disable-gpu-sandbox",
			"--disable-accelerated-2d-canvas",
			"--disable-accelerated-jpeg-decoding",
			"--test-type=ui",
		},
	}
	l.caps.AddChrome(chromeCaps)
}

func (l *Loader) Process(c *cli.Context, domains []string, url string) ([]map[string]string, error) {
	var opts []selenium.ServiceOption
	// domainSources := make([]map[string]string, len(domains))
	var domainSources []map[string]string

	w := c.Writer()
	service, err := selenium.NewChromeDriverService("chromedriver", port, opts...)
	if err != nil {
		cli.Printf(w, "Error starting the ChromeDriver server: %v\n", err)
		return domainSources, err
	}
	defer service.Stop()

	wd, err := selenium.NewRemote(l.caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		cli.Println(w, err)
		return domainSources, err
	}
	defer wd.Quit()

	err = wd.Get(instrumentURL)
	if err != nil {
		cli.Printf(w, "Failed to load page: %s\n", err)
		return domainSources, err
	}

	elem, err := wd.FindElement(selenium.ByID, "nav5")
	if err != nil {
		panic(err)
	}

	elem.Click()

	urlInput, err := wd.FindElement(selenium.ByID, "url")
	if err != nil {
		panic(err)
	}

	for _, domain := range domains {
		domainData := make(map[string]string)
		domainData["domain"] = domain
		urlInput.Clear()
		urlInput.SendKeys(domain + url)
		checkBtn, err := wd.FindElement(selenium.ByID, "su")
		if err != nil {
			panic(err)
		}

		checkBtn.Click()
		cli.Printf(w, "Please wait for %ds to load result ...\n", timeToWaitLoading)
		time.Sleep(timeToWaitLoading * time.Second)

		var frameHtml string
		cli.Printf(w, "[%s]Loading page source ...\n", domain)
		time.Sleep(1 * time.Second)
		frameHtml, err = wd.PageSource()

		if err != nil {
			domainData["err"] = fmt.Sprintf("err: %s", err)
			domainSources = append(domainSources, domainData)
			continue
		}

		instrSlice, err := Parse(c, frameHtml)
		bytes, err := json.MarshalIndent(instrSlice, "", "\t")
		// fmt.Printf("Domain[%s]:\n", domain)
		// fmt.Printf("%s\n", bytes)
		// for _, ist := range instrSlice {
		// 	fmt.Printf("%#v\n", ist)
		// }
		domainData["source"] = fmt.Sprintf("%s", bytes)
		domainData["err"] = ""

		domainSources = append(domainSources, domainData)
	}
	return domainSources, nil
}
