package main

import (
	"DiniSQL/MiniSQL/src/API"
	"DiniSQL/MiniSQL/src/BufferManager"
	"DiniSQL/MiniSQL/src/CatalogManager"
	"DiniSQL/MiniSQL/src/Interpreter/parser"
	"DiniSQL/MiniSQL/src/Interpreter/types"
	"DiniSQL/MiniSQL/src/RecordManager"
	"DiniSQL/Region"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/peterh/liner"
)

const historyCommmandFile = "~/.DiniSQL_history"

//FlushALl 结束时做的所有工作
func FlushALl() {
	BufferManager.BlockFlushAll()                                             //缓存block
	RecordManager.FlushFreeList()                                             //free list写回
	CatalogManager.FlushDatabaseMeta(CatalogManager.UsingDatabase.DatabaseId) //刷新记录长度和余量
}

func InitDB() error {
	err := CatalogManager.LoadDbMeta()
	if err != nil {
		return err
	}
	BufferManager.InitBuffer()

	return nil
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		parts := strings.SplitN(path, "/", 2)
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, parts[1]), nil
	}
	return path, nil
}

func loadHistoryCommand() (*os.File, error) {
	var file *os.File
	path, err := expandPath(historyCommmandFile)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
	}
	return file, err
}

func runShell(r chan<- error) {
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	file, err := loadHistoryCommand()
	if err != nil {
		panic(err)
	}
	defer func() {
		_, err := line.WriteHistory(file)
		if err != nil {
			panic(err)
		}
		_ = file.Close()
	}()
	s := bufio.NewScanner(file)
	for s.Scan() {
		//fmt.Println(s.Text())
		line.AppendHistory(s.Text())
	}
	InitDB()

	StatementChannel := make(chan types.DStatements, 500)  //用于传输操作指令通道
	FinishChannel := make(chan string, 500)                //用于api执行完成反馈通道
	FlushChannel := make(chan struct{})                    //用于每条指令结束后协程flush
	go API.HandleOneParse(StatementChannel, FinishChannel) //begin the runtime for exec
	go BufferManager.BeginBlockFlush(FlushChannel)
	var beginSQLParse = false
	var sqlText = make([]byte, 0, 100)
	for { //each sql
	LOOP:
		beginSQLParse = false
		sqlText = sqlText[:0]
		var input string
		var err error
		for { //each line
			if beginSQLParse {
				input, err = line.Prompt("      -> ")
			} else {
				input, err = line.Prompt("DiniSQL> ")
			}
			if err != nil {
				if err == liner.ErrPromptAborted {
					goto LOOP
				}
			}
			trimInput := strings.TrimSpace(input) //get the input without front and backend space
			if len(trimInput) != 0 {
				line.AppendHistory(input)
				if !beginSQLParse && (trimInput == "quit" || strings.HasPrefix(trimInput, "quit;")) {
					close(StatementChannel)
					for _ = range FinishChannel {

					}
					FlushALl()
					r <- err
					return
				}
				sqlText = append(sqlText, append([]byte{' '}, []byte(trimInput)[0:]...)...)
				if !beginSQLParse {
					beginSQLParse = true
				}
				if strings.Contains(trimInput, ";") {
					break
				}
			}
		}
		beginTime := time.Now()
		err = parser.Parse(strings.NewReader(string(sqlText)), StatementChannel)
		//fmt.Println(string(sqlText))
		if err != nil {
			fmt.Println(err)
			continue
		}
		result := <-FinishChannel //等待指令执行完成
		fmt.Println(result)
		durationTime := time.Since(beginTime)
		fmt.Println("Finish operation at: ", durationTime)
		FlushChannel <- struct{}{} //开始刷新cache
	}
}

// func main() {
// 	//errChan 用于接收shell返回的err
<<<<<<< HEAD
// 	Region.InitRegionServer()
=======
>>>>>>> 45d7f851106a218e0eaff95b5b10c1ad95fde9cb
// 	errChan := make(chan error)
// 	go runShell(errChan) //开启shell协程
// 	err := <-errChan
// 	fmt.Println("bye")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }
<<<<<<< HEAD
=======

func main() {
	var endpoints = []string{"127.0.0.1:2379"}
	ser, err := Region.NewServiceRegister(endpoints, "/web/node1", "localhost:8000", 5)
	if err != nil {
		log.Fatalln(err)
	}
	//监听续租相应chan
	go ser.ListenLeaseRespChan()
	select {
	// case <-time.After(20 * time.Second):
	// 	ser.Close()
	}
}
>>>>>>> 45d7f851106a218e0eaff95b5b10c1ad95fde9cb
