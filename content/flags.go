package content

import "turbojet/cli"

const (
	NameFlagName            = "content-name"
	FromFlagName            = "content-from"
	ContentIDFlagName       = "content-id"
	ContentTypeFlagName     = "content-type"
	ContentPropertyFlagName = "content-property"
)

func AddFlags(fs *cli.FlagSet) {
	fs.Add(NewNameFlag())
	fs.Add(NewFromFlag())
	fs.Add(NewContentIDFlag())
	fs.Add(NewContentTypeFlag())
	fs.Add(NewContentPropertyFlag())
}

func NewNameFlag() *cli.Flag {
	return &cli.Flag{
		Category:   "content",
		Name:       NameFlagName,
		Persistent: true,
		Short:      "use `--content-name <contentName> to select Context Package`",
	}
}

func NameFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(NameFlagName)
}

func NewFromFlag() *cli.Flag {
	return &cli.Flag{
		Category:   "content",
		Name:       FromFlagName,
		Persistent: true,
		Short:      "use `--content-from <directory> to load content files",
	}
}

func FromFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(FromFlagName)
}

func NewContentIDFlag() *cli.Flag {
	return &cli.Flag{
		Category:   "content",
		Name:       ContentIDFlagName,
		Persistent: true,
		Short:      "use `--content-id <contentID> to load content",
	}
}

func ContentIDFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(ContentIDFlagName)
}

func NewContentTypeFlag() *cli.Flag {
	return &cli.Flag{
		Category:   "content",
		Name:       ContentTypeFlagName,
		Persistent: true,
		Short:      "use `--content-type <mini_lobby, lobby or game> to load type content",
	}
}

func ContentTypeFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(ContentTypeFlagName)
}

func NewContentPropertyFlag() *cli.Flag {
	return &cli.Flag{
		Category:   "content",
		Name:       ContentPropertyFlagName,
		Persistent: true,
		Short:      "use `--content-property <propertyID> to load content",
	}
}

func ContentPropertyFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(ContentPropertyFlagName)
}
