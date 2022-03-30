package transport

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ntpusdk/setup"
)

/*
	Буфер приема/ответа запроса на точное время
	позиция	размер 	комментарий
	0			10	Идентификатор устройства
	10			1	число
	11			1   месяц
	12			1	год дву мл цифры от 2000+
	13 			1	час
	14			1   минута
	15			1 	секунда
*/
func workerNtp(socket net.Conn) {
	defer socket.Close()
	for {
		buffer := make([]byte, 16)
		socket.SetReadDeadline(time.Now().Add(time.Minute))
		size, err := socket.Read(buffer)
		if err != nil {
			logger.Error.Printf("Ntp user %s %s", socket.RemoteAddr().String(), err.Error())
			return
		}
		if size != 16 {
			logger.Error.Printf("Ntp user %s bad lenght %d", socket.RemoteAddr().String(), size)
			return
		}
		// Читаем id устройства
		id := make([]byte, 8)
		for i := 0; i < len(id); i++ {
			id[i] = buffer[i+2] + '0'
		}
		lid, err := strconv.Atoi(string(id))
		if err != nil {
			logger.Error.Printf("Ntp user %s bad id", socket.RemoteAddr().String())
			return
		}
		dtime := TakeDate(buffer, 10)
		if !EqualTime(dtime, time.Now()) {
			PutDate(time.Now(), buffer, 10)
			logger.Info.Printf("Устройство %d его время %s устанавливаем %s", lid, dtime.Format(time.RFC3339), time.Now().Format(time.RFC3339))
		} else {
			logger.Info.Printf("Устройство %d его время корректно", lid)
		}
		socket.SetWriteDeadline(time.Now().Add(time.Second))
		size, err = socket.Write(buffer)
		if err != nil {
			logger.Error.Printf("Ntp user id %d %s %s", lid, socket.RemoteAddr().String(), err.Error())
			return
		}
		if size != len(buffer) {
			logger.Error.Printf("Ntp user %d %s recive bad lenght %d", lid, socket.RemoteAddr().String(), size)
			return
		}
	}
}
func EqualTime(ti, tj time.Time) bool {
	if ti.Year() != tj.Year() {
		return false
	}
	if ti.Month() != tj.Month() {
		return false
	}
	if ti.Day() != tj.Day() {
		return false
	}
	if ti.Hour() != tj.Hour() {
		return false
	}
	if ti.Minute() != tj.Minute() {
		return false
	}
	if ti.Second() != tj.Second() {
		return false
	}
	return true
}
func TakeDate(buffer []byte, pos int) time.Time {
	year := int(buffer[pos+2]) + 2000
	month := time.Month(int(buffer[pos+1]))
	day := int(buffer[pos])
	hour := int(buffer[pos+3])
	minut := int(buffer[pos+4])
	sec := int(buffer[pos+5])
	location, _ := time.LoadLocation("Local")
	return time.Date(year, month, day, hour, minut, sec, 0, location)
}
func PutDate(t time.Time, buffer []byte, pos int) {
	year, month, day := t.Date()
	hour := t.Hour()
	min := t.Minute()
	sec := t.Second()
	buffer[pos] = uint8(day)
	buffer[pos+1] = uint8(month)
	buffer[pos+2] = uint8(year % 100)
	buffer[pos+3] = uint8(hour)
	buffer[pos+4] = uint8(min)
	buffer[pos+5] = uint8(sec)
}

func ListenExternalDevices() {
	ln, err := net.Listen("tcp4", fmt.Sprintf(":%d", setup.Set.NtpPort))
	if err != nil {
		logger.Error.Printf("Ошибка открытия ntp port %s", err.Error())
		return
	}
	for {
		socket, err := ln.Accept()
		if err != nil {
			logger.Error.Printf("Accept %s", err.Error())
			continue
		}
		logger.Info.Printf("Новый запрос на корректировку времени %s", socket.RemoteAddr().String())
		go workerNtp(socket)
	}
}
