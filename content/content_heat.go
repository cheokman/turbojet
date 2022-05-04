package content

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"time"
	"turbojet/cli"
	"turbojet/heater"
	"turbojet/platform"
	"turbojet/provider"
)

const (
	defaultIContentURL = "https://game.cache.com"
)

func NewContentHeatCommand() *cli.Command {
	return &cli.Command{
		Name:  "heat",
		Short: "heat content",
		Usage: "heat --content-id <contentID> --content-property <propertyID>",
		Run: func(c *cli.Context, args []string) error {
			doContentHeat(c, args)
			return nil
		},
	}
}

func doContentHeat(c *cli.Context, args []string) {
	w := c.Writer()
	propertyID, ok := ContentPropertyFlag(c.Flags()).GetValue()
	if !ok {
		cli.Printf(w, "Need to specify a property to heat \n")
		return
	}
	contentID, ok := ContentIDFlag(c.Flags()).GetValue()

	if !ok {
		cList, _, ok1 := LoadContentsFromStorage(c, GetStoragePath()+string(os.PathSeparator)+storageFile)
		if ok1 != nil {
			cli.Printf(w, "Load Content IDs Error!\n")
		}

		// var cIDs []string
		start := time.Now()
		for _, cID := range cList {
			// cIDs = append(cIDs, cID.ID)
			cID.Heat(c, propertyID)
		}
		processTime := time.Since(start)
		// cli.Printf(w, "To be Developed for heat all %#v\n", cIDs)
		cli.Printf(w, "Processed in total %s\n", processTime)
		return
	} else {
		contentType, ok := ContentTypeFlag(c.Flags()).GetValue()
		if !ok {
			contentType = "game"
		}
		cli.Printf(w, "Heating content ID %s and type %s\n", contentID, contentType)
		content, _, ok := LoadContentFromStorageByID(c, GetStoragePath()+string(os.PathSeparator)+storageFile, contentID)
		if !ok {
			cli.Printf(w, "load content error\n")
			return
		}
		content.Heat(c, propertyID)

		return
	}
}

func HeatAll(ctx *cli.Context, content []Content) error {
	fmt.Printf("In HeatAll")
	return nil
}

func (c *Content) Heat(ctx *cli.Context, propertyID string) error {
	w := ctx.Writer()
	// profile, err := config.LoadProfileWithContext(ctx)
	// if err != nil {
	// 	cli.Printf(w, "load profile error: %s\n", err)
	// }
	var ips []string
	var domains []string
	cNode := make(map[string][]string)
	var perr error
	var rURL string
	switch c.ContentType {
	case platformLobby:
		rURL, perr = platform.GetLobbyCDNPath(ctx, defaultIContentURL, propertyID)
	case miniLobby:
		rURL, perr = platform.GetMiniLobbyCDNPath(ctx, defaultIContentURL, propertyID)
	case game:
		rURL, perr = platform.GetGameCDNPath(ctx, defaultIContentURL, propertyID, c.ID)
		fmt.Printf("Game: %s\n", rURL)
		fmt.Printf("error: %s\n", perr)
	}
	if perr != nil {
		cli.Printf(w, "content[%s] of property[%s] get CDN Path error: %s\n", c.ID, propertyID, perr)
		return nil
	}

	u, err := url.Parse(rURL)
	if err != nil {
		cli.Printf(w, "content[%s] of property[%s] parse CDN Path error: %s\n", c.ID, propertyID, err)
	}

	pathSl := strings.Split(u.Path, "/")
	cntPath := ""
	pathSlLen := len(pathSl)
	if pathSlLen > 0 {
		cntPath = strings.Join(pathSl[1:pathSlLen-1], "/")
	}

	provider, _, _ := provider.LoadProviderWithContext(ctx)
	for _, d := range provider.Domains {

		cli.Printf(w, "Domains: %s\n", d.Name)
		instruments, _ := provider.LoadSource(ctx, d.Name)
		currentIPs := instruments.GetIPs()
		ips = append(ips, currentIPs...)
		cNode[d.Name] = currentIPs
		domains = append(domains, d.Name)
	}

	cli.Printf(w, "Property ID: %s\n", propertyID)
	cli.Printf(w, "Content ID: %s\n", c.ID)
	cli.Printf(w, "URL: %s\n", rURL)
	h := heater.NewHeater(cNode, cntPath, c.GetRelativePaths(), c.ContentType)
	jobs := h.Process(ctx)
	tw := tabwriter.NewWriter(w, 8, 0, 1, ' ', 0)
	fmt.Fprint(tw, "\nIP\t| Total DLs\t| Total DL Suc\t| Total DL Err\t| Total DL Size\t| Total Prc Time \t| Max Prc Time \t| Avg Prc Time\t| Min Prc Time\t|\n")
	fmt.Fprintf(tw, "-------\t| --------\t| --------\t| --------\t| --------\t|  --------\t| --------\t| --------\t| --------\n")
	for _, j := range jobs {
		js := j.GetSummary()
		ips := js.IPSummaries
		for _, ip := range ips {
			fmt.Fprintf(tw, "%s\t| %d\t| %d\t| %d\t| %d\t| %s\t| %s\t| %s\t| %s\n", ip.IP, ip.TotalDownload, ip.TotalDownloadSuc, ip.TotalDownloadErr, ip.TotalDownloadSize, ip.TotalProcessTime, ip.MaxProcessTime, ip.AvgProcessTime, ip.MinProcessTime)
		}
	}
	tw.Flush()
	return nil
}
