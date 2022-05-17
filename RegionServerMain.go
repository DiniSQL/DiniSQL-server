package main

import (
	"DiniSQL/MiniSQL/src/BufferManager"
	"DiniSQL/MiniSQL/src/CatalogManager"
	"DiniSQL/MiniSQL/src/RecordManager"
	Region "DiniSQL/Region"
)

func InitDB() error {
	err := CatalogManager.LoadDbMeta()
	if err != nil {
		return err
	}
	BufferManager.InitBuffer()

	return nil
}

func FlushALl() {
	BufferManager.BlockFlushAll()                                             //缓存block
	RecordManager.FlushFreeList()                                             //free list写回
	CatalogManager.FlushDatabaseMeta(CatalogManager.UsingDatabase.DatabaseId) //刷新记录长度和余量
}

func main() {
	InitDB()
	Region.InitRegionServer()
}
