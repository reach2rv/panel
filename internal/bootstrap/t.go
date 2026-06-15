package bootstrap

import (
	"os"

	"github.com/leonelquinteros/gotext"

	"github.com/acepanel/panel/v3/pkg/config"
	"github.com/acepanel/panel/v3/pkg/embed"
)

func NewT(conf *config.Config) (*gotext.Locale, error) {
	if conf.App.Locale == "en" {
		_ = os.Setenv("LANG", "en_US.UTF-8")
		_ = os.Setenv("LC_ALL", "en_US.UTF-8")
		_ = os.Setenv("LC_MESSAGES", "en_US.UTF-8")
	}

	l := gotext.NewLocaleFSWithPath(conf.App.Locale, embed.LocalesFS, "locales")
	l.AddDomain("backend")

	return l, nil
}
