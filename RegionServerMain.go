package main

import (
	"DiniSQL/MiniSQL/src/BufferManager"
	"DiniSQL/MiniSQL/src/CatalogManager"
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

func main() {
	InitDB()
	Region.InitRegionServer()
}
