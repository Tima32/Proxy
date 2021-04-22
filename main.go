package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"unsafe"
)

func main() {
	fmt.Println("Hello, I'm proxy server")

	go commandExit()

	listener, _ := net.Listen("tcp", ":3202") // открываем
	if listener == nil {
		panic("Error open tcp server")
	}

	for {
		conn, err := listener.Accept() // принимаем TCP-соединение от клиента и создаем новый сокет
		if err != nil {
			continue
		}
		go handleClient(conn) // обрабатываем запросы клиента в отдельной го-рутине
	}
}
func commandExit() {
	for {
		var s string
		fmt.Scan(&s)
		if s == "exit" || s == "stop" {
			os.Exit(0)
		}
	}
}

//session - сессия
type session struct {
	conn       net.Conn //сокет
	serverConn net.Conn //сокет с сервером
	SucretKey  []byte   //ключ шифрования
	Crypter    CrypterString
	Uncrypter  CrypterString
}

//создание общего ключа
func deffieHelman(conn net.Conn) ([]byte, error) {
	var key32 [128]uint32
	key := (*[512]byte)(unsafe.Pointer(&key32))[:]

	var publicKeyClient [128]uint64
	var privateKey [128]uint64
	var publicKey [128]uint64

	//получение публичного ключа клиента
	err := binary.Read(conn, binary.LittleEndian, &publicKeyClient)
	if err != nil {
		return key, err
	}

	//генерация приватного и публичного ключа
	for i := 0; i < 128; i++ {
		privateKey[i], publicKey[i] = GeneratePublicPrivateKey()
	}

	//тправка публичного ключа клиенту
	err = binary.Write(conn, binary.LittleEndian, publicKey)
	if err != nil {
		return key, err
	}

	//вычисление ключа
	for i := 0; i < 128; i++ {
		key32[i] = FindSucretKey(publicKeyClient[i], privateKey[i])
	}

	return key, nil
}

//обработчик нового клиента
func handleClient(conn net.Conn) {
	var Session session
	Session.conn = conn

	//получение ключа шифрования
	sucretKey, err := deffieHelman(Session.conn)
	if err != nil {
		fmt.Println("deffieHelman err: ", err)
		conn.Close()
		return
	}
	Session.SucretKey = sucretKey
	//инициализация криптографии
	Session.Crypter = NewCrypterString()
	Session.Uncrypter = NewCrypterString()
	Session.Crypter.SetKey(sucretKey)
	Session.Uncrypter.SetKey(sucretKey)
	fmt.Println("Sucret key: ", sucretKey[0], " ", sucretKey[511])

	connS, errc := net.Dial("tcp", "127.0.0.1:32111")
	if errc != nil {
		fmt.Println("Server connect err: ", errc)
		conn.Close()
		Session.conn.Close()
		return
	}
	Session.serverConn = connS
	go sToC(Session)

	//client to server
	rBuff := make([]byte, 1048576)
	for {
		nr, err := conn.Read(rBuff)
		if err != nil {
			fmt.Println("Error c to s read")
			conn.Close()
			Session.conn.Close()
			return
		}
		wBuff := rBuff[0:nr]
		wBuff = Session.Uncrypter.CryptBytes(wBuff)

		_, errs := Session.serverConn.Write(wBuff)
		if errs != nil {
			fmt.Println("Error c to s write")
			conn.Close()
			Session.conn.Close()
			return
		}
	}
}

func sToC(Session session) {

	rBuff := make([]byte, 1048576)
	for {
		nr, err := Session.serverConn.Read(rBuff)
		if err != nil {
			fmt.Println("Error s to c read")
			Session.serverConn.Close()
			Session.conn.Close()
			return
		}
		wBuff := rBuff[0:nr]
		wBuff = Session.Crypter.CryptBytes(wBuff)

		_, errs := Session.conn.Write(wBuff)
		if errs != nil {
			fmt.Println("Error s to c write")
			Session.serverConn.Close()
			Session.conn.Close()
			return
		}
	}
}
