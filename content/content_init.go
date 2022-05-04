package content

import (
	"fmt"
	"os"
	"turbojet/cli"
)

const (
	defaultContentCacheDirectory = "content"
)

func NewContentInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Short: "init content throught cache or API",
		Usage: "init --content-from <directory> --cache ...",
		Run: func(c *cli.Context, args []string) error {
			contentFrom, _ := FromFlag(c.Flags()).GetValue()
			doContentInit(c, contentFrom)
			return nil
		},
	}
}

func doContentInit(c *cli.Context, cntFrom string) {
	w := c.Writer()
	storage, err := LoadStorage(GetStoragePath()+string(os.PathSeparator)+storageFile, w)
	if err != nil {
		cli.Printf(w, "Content init loading storage error: %s\n", err)
	}

	fmt.Printf("Content from: %s\n", cntFrom)
	contents, err := LoadContentFromLocal(c, cntFrom)
	if err != nil {
		cli.Printf(w, "load content from local['%s'] error: %s", cntFrom, err)
		return
	}

	for _, cnt := range contents {
		storage.PutContent(*cnt)
		// for _, cf := range cnt.ContentFiles {
		// 	fmt.Printf("%#v\n", cf)
		// }
	}
	SaveStorage(storage)
	return
}
