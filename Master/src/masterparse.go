package main

import (
	"DiniSQL/MiniSQL/src/Interpreter/types"
	"DiniSQL/MiniSQL/src/Utils/Error"
	"fmt"
)

//HandleOneParse 用来处理parse处理完的DStatement类型  dataChannel是接收Statement的通道,整个mysql运行过程中不会关闭，但是quit后就会关闭
//stopChannel 用来发送同步信号，每次处理完一个后就发送一个信号用来同步两协程，主协程需要接收到stopChannel的发送后才能继续下一条指令，当dataChannel
//关闭后，stopChannel才会关闭
func Parse2Statement(dataChannel <-chan types.DStatements, stopChannel chan<- string, tableChannel chan<- []string, stmtChannel chan<- types.DStatements) {
	var _ Error.Error
	for statement := range dataChannel {
		//fmt.Println("statement in Parse2Statement:", statement)
		var ret string
		stmtChannel <- statement
		switch statement.GetOperationType() {
		case types.CreateDatabase:
			fmt.Println("create database")
		case types.UseDatabase:
			fmt.Println("use database")
		case types.CreateTable:
			fmt.Println("create table")
			fmt.Println(statement.(types.CreateTableStatement).TableName)
			tableChannel <- []string{statement.(types.CreateTableStatement).TableName}

		case types.CreateIndex:
			tableChannel <- []string{statement.(types.CreateIndexStatement).TableName}

		case types.DropTable:
			tableChannel <- []string{statement.(types.DropTableStatement).TableName}

		case types.DropIndex:
			tableChannel <- []string{statement.(types.DropIndexStatement).TableName}

		case types.DropDatabase:

		case types.Insert:
			tableChannel <- []string{statement.(types.InsertStament).TableName}

		case types.Update:
			tableChannel <- []string{statement.(types.UpdateStament).TableName}

		case types.Delete:
			tableChannel <- []string{statement.(types.DeleteStatement).TableName}

		case types.Select:
			fmt.Println("table: ", statement.(types.SelectStatement).TableNames[0])
			tableChannel <- statement.(types.SelectStatement).TableNames

		case types.ExecFile:
			fmt.Println("execfile")
			//fmt.Println(err)
			stopChannel <- ret
		}
		//close(stopChannel)
	}
}
