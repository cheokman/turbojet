package content

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"turbojet/cli"
)

func NewContentListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Short: "list content",
		Usage: "list",
		Run: func(ctx *cli.Context, args []string) error {
			if len(args) > 0 {
				return cli.NewInvalidCommandError(args[0], ctx)
			}
			doContentList(ctx)
			return nil
		},
	}
}

func doContentList(ctx *cli.Context) {
	w := ctx.Writer()
	loadPath := GetStoragePath() + string(os.PathSeparator) + storageFile
	storage, err := LoadStorage(loadPath, w)
	if err != nil {
		cli.Printf(w, "Content init loading storage error: %s\n", err)
	}
	fmt.Printf("Loaded Content From %s\n", loadPath)
	tw := tabwriter.NewWriter(w, 8, 0, 1, ' ', 0)
	fmt.Fprint(tw, "\nType\t| ID\t| Total Files\t| Total Size\t| Types Distribution\t| Size Distributions\n")
	fmt.Fprintf(tw, "-------\t| --------\t| --------\t| --------\t| --------\t| -------\n")
	for _, cnt := range storage.Contents {
		cnt.LoadContent(ctx)
		typeDis := cnt.GetTypeDistribution()
		tdStr := ""
		for k, v := range typeDis {
			tdStr = tdStr + fmt.Sprintf(" %s(%d)", strings.Replace(k, ".", "", -1), v)
		}
		sizeDis := cnt.GetSizeDistribution()
		sdStr := ""
		for k, v := range sizeDis {
			sdStr = sdStr + fmt.Sprintf(" %s(%d)", strings.Replace(k, ".", "", -1), v)
		}
		fmt.Fprintf(tw, "%s\t|%s\t|%d\t|%d\t|%s\t|%s\t\n", cnt.ContentType, cnt.ID, cnt.GetTotalFiles(), cnt.GetTotalSize(), tdStr, sdStr)
	}
	tw.Flush()
}
