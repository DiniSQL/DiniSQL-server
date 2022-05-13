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

type simpleRegion struct {
	visitCount int
	regionID   int
}

func (r *region) simplify() simpleRegion {
	s := simpleRegion{visitCount: r.visitCount, regionID: r.regionID}
	return s
}
