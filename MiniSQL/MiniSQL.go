package MiniSQL

import (
	"DiniSQL/MiniSQL/src/API"
	"DiniSQL/MiniSQL/src/BufferManager"
	"DiniSQL/MiniSQL/src/CatalogManager"
	"DiniSQL/MiniSQL/src/Interpreter/parser"
	"DiniSQL/MiniSQL/src/Interpreter/types"
	"DiniSQL/MiniSQL/src/RecordManager"
	"bufio"
	"fmt"
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

func RunShell(r chan<- error) {
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
		fmt.Println(result[1:])
		durationTime := time.Since(beginTime)
		fmt.Println("Finish operation at: ", durationTime)
		FlushChannel <- struct{}{} //开始刷新cache
	}
}
