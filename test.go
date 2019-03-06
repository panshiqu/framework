package main

import (
	"./define"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

func main()  {
	conn,err := net.Dial("tcp4", "127.0.0.1:10030")
	if err != nil {
		fmt.Println("dial failed:"+err.Error())
	}
	//network.SendMessage()
	defer conn.Close()
	register := define.FastRegister{
		Account:"king123",
		Gender:  1,
		Icon: 1,
		Password: "123456",
		Name: "big king2",
		IP: "127.0.0.1",
	}

	data,_ := json.Marshal(register)
	fmt.Println(string(data))

	msg := getMessage(define.LoginCommon,define.LoginFastRegister, []byte(data))
	res, err := conn.Write(msg)
	fmt.Println("res:", res)
	return
}


func getMessage(mcmd uint16, scmd uint16, data []byte) []byte {
	size := len(data) + 6
	message := make([]byte, size)
	binary.BigEndian.PutUint16(message, uint16(size))
	binary.BigEndian.PutUint16(message[2:], mcmd)
	binary.BigEndian.PutUint16(message[4:], scmd)
	copy(message[6:], data)

	return message
}