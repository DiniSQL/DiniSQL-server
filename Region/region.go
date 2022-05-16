package Region

import (
	"DiniSQL/MiniSQL/src/CatalogManager"
	//"DiniSQL/MiniSQL/src/BufferManager"
)

type region struct {
	tableInfo  CatalogManager.TableCatalog
	visitCount int
	regionID   int
}
