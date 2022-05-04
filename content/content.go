package content

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"turbojet/cli"
)

const (
	latestVersion = "latest"
	fromLocal     = "local"
	fromURL       = "url"
	miniLobby     = "miniLobby"
	platformLobby = "platformLobby"
	game          = "game"
)

type Content struct {
	Name         string         `json:"name"`
	ID           string         `json:"id"`
	Version      string         `json:"version"`
	Path         string         `json:"path"`
	File         string         `json:"file"`
	From         string         `json:"from"`
	ContentType  string         `json:"content_type"`
	ContentFiles []*ContentFile `json:"content_files"`
}

func NewContent(path string, from string) *Content {
	if from == fromLocal {
		id := parseContentFileID(path)
		return &Content{
			Name:        parseContentFileName(path),
			ID:          id,
			ContentType: GetType(id),
			File:        filepath.Base(path),
			From:        fromLocal,
			Version:     latestVersion,
			Path:        path,
		}
	} else {
		return &Content{}
	}
}

func GetType(id string) string {
	switch id {
	case "9006":
		return platformLobby
	case "9007":
		return miniLobby
	default:
		return game
	}
}

func NewContentCommand() *cli.Command {
	c := &cli.Command{
		Name:  "content",
		Short: "content information and operations",
		Usage: "content --content-name <Content Name> --content-type <mini_lobby/lobby/game>",
		Run: func(c *cli.Context, args []string) error {
			contentName, _ := NameFlag(c.Flags()).GetValue()
			return doContent(c, contentName)
		},
	}
	c.AddSubCommand(NewContentInitCommand())
	c.AddSubCommand(NewContentHeatCommand())
	c.AddSubCommand(NewContentListCommand())
	return c
}

func doContent(c *cli.Context, contentName string) error {
	fmt.Printf("Content name: %s\n", contentName)
	return nil
}

func LoadContentFromLocal(c *cli.Context, folder string) ([]*Content, error) {
	var ctnSlice []*Content
	foldInfo, err := os.Stat(folder)
	if os.IsNotExist(err) {
		return ctnSlice, err
	}

	if !foldInfo.IsDir() {
		return ctnSlice, NewInvalidCacheFromError(folder)
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return ctnSlice, NewInvalidCacheFromError(folder)
	}

	for _, f := range files {
		fn := folder + "/" + f.Name()
		ctn, err := loadContentFromFile(c, fn)
		if err != nil {
			cli.Printf(c.Writer(), "load content from file[%s] error: %s", fn, err)
		}
		ctnSlice = append(ctnSlice, ctn)
	}
	return ctnSlice, nil
}

func loadContentFromFile(ctx *cli.Context, file string) (*Content, error) {
	content := NewContent(file, fromLocal)
	fmt.Printf("content: %#v\n\n", content)
	err := content.LoadContent(ctx)
	if err != nil {
		return content, err
	}
	return content, nil
}

func parseContentFileName(f string) string {
	fileSplited := strings.Split(filepath.Base(f), ".")
	if len(fileSplited) > 0 {
		return fileSplited[0]
	}
	return ""
}

func parseContentFileID(f string) string {
	fileSplited := strings.Split(filepath.Base(f), ".")
	if len(fileSplited) > 0 {
		return fileSplited[0]
	}
	return ""
}

func (c *Content) GetRelativePaths() []string {
	var paths []string
	for _, cf := range c.ContentFiles {
		var path string
		if cf.Path == "." {
			path = cf.Name
		} else {
			path = fmt.Sprintf("%s/%s", cf.Path, cf.Name)
		}
		paths = append(paths, path)
	}
	return paths
}

func (c *Content) GetTotalFiles() int {
	return len(c.ContentFiles)
}

func (c *Content) GetTotalSize() int {
	var totalSize int
	for _, cf := range c.ContentFiles {
		totalSize = totalSize + cf.Size
	}
	return totalSize
}

func (c *Content) GetTypeDistribution() map[string]int {
	typeDis := make(map[string]int)
	for _, cf := range c.ContentFiles {
		ext := cf.Ext
		typeDis[ext]++
	}
	return typeDis
}

func (c *Content) GetSizeDistribution() map[string]int {
	sizeDis := make(map[string]int)
	for _, cf := range c.ContentFiles {
		ext := cf.Ext
		size := cf.Size
		sizeDis[ext] += size
	}
	return sizeDis
}

func (c *Content) LoadContent(ctx *cli.Context) error {
	if _, err := os.Stat(c.Path); os.IsNotExist(err) {
		cli.Printf(ctx.Writer(), "Create content from file['%s'] error: %s\n", err)
	}

	bytes, err := ioutil.ReadFile(c.Path)
	var assetFiles []string
	err = json.Unmarshal(bytes, &assetFiles)
	if err != nil {
		return err
	}
	for _, f := range assetFiles {
		cf, err := NewContentFile(f)
		if err != nil {
			cli.Printf(ctx.Writer(), "read content file['%s'] error: %s\n", c.File, err)
			continue
		}
		c.ContentFiles = append(c.ContentFiles, cf)
	}
	return nil
	// file, err := os.Open(c.Path)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	cf, err := NewContentFile(scanner.Text())
	// 	if err != nil {
	// 		cli.Printf(ctx.Writer(), "read content file['%s'] error: %s\n", c.File, err)
	// 		continue
	// 	}
	// 	c.ContentFiles = append(c.ContentFiles, cf)
	// }
	// return nil
}
