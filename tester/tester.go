package tester

import (
	"fmt"
	"net"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ntpusdk/setup"
	"github.com/ruraomsk/ntpusdk/transport"
)

func RunTester() {
	for {
		con := fmt.Sprintf("127.0.0.1:%d", setup.Set.NtpPort)
		socket, err := net.Dial("tcp", con)
		if err != nil {
			logger.Error.Printf("connect %s %s ", con, err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
	loop:
		for {
			time.Sleep(20 * time.Second)
			buffer := make([]byte, 16)
			buffer[5] = 2
			transport.PutDate(time.Now(), buffer, 10)
			socket.SetWriteDeadline(time.Now().Add(time.Second))
			size, err := socket.Write(buffer)
			if err != nil {
				logger.Error.Printf("Ошибка передачи %s ", err.Error())
				break loop
			}
			if size != 16 {
				logger.Error.Printf("Ошибка передачи оправлено %d ", size)
				break loop
			}
			socket.SetReadDeadline(time.Now().Add(time.Second))
			size, err = socket.Read(buffer)
			if err != nil {
				logger.Error.Printf("Ошибка приема %s ", err.Error())
				break loop
			}
			if size != 16 {
				logger.Error.Printf("Ошибка приема получено %d ", size)
				break loop
			}
			logger.Info.Printf("Время на сервере %s", transport.TakeDate(buffer, 10).Format(time.RFC3339))

		}
		socket.Close()
	}
}
