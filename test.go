package main

import (
	"./define"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/panshiqu/framework/network"
	"net"
)

func main()  {
	conn,err := net.Dial("tcp4", "127.0.0.1:10030")
	if err != nil {
		fmt.Println("dial failed:"+err.Error())
	}
	//network.SendMessage()
	defer conn.Close()

	registerCheck(conn)
	//
	//register := define.FastRegister{
	//	Account:"king111",
	//	Gender:  1,
	//	Icon: 1,
	//	Password: "123456",
	//	Name: "big king",
	//	IP: "127.0.0.1",
	//}
	//doSendMessage(conn,define.LoginCommon,define.LoginFastRegister,register)

	return
}

func doSendMessage(conn net.Conn,mcmd uint16, scmd uint16,inMsg interface{}) {
	data,_ := json.Marshal(inMsg)
	fmt.Println(string(data))
	msg := getSendMessage(mcmd,scmd, []byte(data))
	res, _ := conn.Write(msg)
	fmt.Println("res:", res)
	_,_,buf,_ := network.RecvMessage(conn)
	fmt.Println(string(buf))
}

func registerCheck(conn net.Conn) {
	check := define.FastRegisterCheck{
		Account:"king",
		Name:"wong",
	}
	doSendMessage(conn,define.LoginCommon,define.LoginRegisterCheck, check)
}

func getSendMessage(mcmd uint16, scmd uint16, data []byte) []byte {
	size := len(data) + 6
	message := make([]byte, size)
	binary.BigEndian.PutUint16(message, uint16(size))
	binary.BigEndian.PutUint16(message[2:], mcmd)
	binary.BigEndian.PutUint16(message[4:], scmd)
	copy(message[6:], data)

	return message
}