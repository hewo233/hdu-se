package Init

import (
	"github.com/hewo233/hdu-se/db"
	"github.com/hewo233/hdu-se/models"
	"github.com/hewo233/hdu-se/shared/consts"
)

func AllInit() {
	db.Init()
	models.SetCozeToken(consts.CozeTokenFile)
}
