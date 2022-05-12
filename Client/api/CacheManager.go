package main

var cache map[string]string
func newCache(){
	cache=make(map[string]string)
}
func flushCache(){
	cache=make(map[string]string)
}
func setCache(table string,server string){
	cache[table]=server
}
func getCache(table string) string{
	// var result,ok string
	result,ok:=cache[table]
	if (ok){
		return result
	}else{
		return ""
	}
}

// import "fmt"
