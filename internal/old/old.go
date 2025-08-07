package old

import (
	"github.com/nettojulio/ufape-crawler-golang/internal/crawler"
	"github.com/nettojulio/ufape-crawler-golang/utils/configs"
)

func old() {
	configs.DisplayConfigs()
	crawler.Execute()
}
