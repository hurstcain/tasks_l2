package main

import (
	"bufio"
	"context"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

// Connect - создает клиентское подключение к веб-сокет серверу.
func Connect(timeout int, url string) *websocket.Conn {
	// Опции для подключения к веб-сокет серверу.
	// Здесь устанавливается таймаут на подключение к серверу (timeout секунд).
	var DefaultDialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: time.Duration(timeout) * time.Second,
	}

	// Контекст, которые завершается по истечению timeout секунд.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	log.Printf("Connecting to %s...\n", url)
	for {
		select {
		// Если спустя timeout секунд не удалось соединиться с сервером, то завершаем программу.
		case <-ctx.Done():
			log.Fatalln("Server connection timeout expired.")
		default:
			// Пробуем подключиться к серверу.
			conn, _, _ := DefaultDialer.Dial(url, nil)
			// Если соединиться не удалось, то переменная conn будет пуста.
			// Тогда переходим на следующую итерацию цикла.
			// И так до тех пор, пока не установится соединение с сервером, либо пока не пройдет таймаут.
			if conn == nil {
				continue
			}
			log.Printf("Connected to %s\n", url)
			return conn
		}
	}
}

func main() {
	// Таймаут на подключение к серверу.
	// По умолчанию - 10 секунд.
	timeout := flag.Int("timeout", 10, "Timeout for connecting to the server.")
	argsLen := len(os.Args)
	// Url хоста, к которому нужно будет подключаться.
	// Состоит из двух последних аргументов командной строки: ip или доменного имени и порта.
	url := "ws://" + os.Args[argsLen-2] + ":" + os.Args[argsLen-1]
	// echo websocket server: ws.ifelse.io:80
	flag.Parse()

	conn := Connect(*timeout, url)
	defer conn.Close()

	// Канал, который закрывается, если сервер закрыл подключение или если клиент решил закончить отправку данных.
	// После того как канал закрывается, программа завершает работу.
	exitCh := make(chan struct{}, 1)

	// Горутина, где принимаются сообщения с сервера.
	go func() {
		for {
			// Принимаем сообщение.
			_, msg, err := conn.ReadMessage()
			// Если произошла ошибка, то закрываем канал exitCh и завершаем работу горутины.
			if err != nil {
				log.Println("Error in receive:", err)
				close(exitCh)
				return
			}
			log.Printf("Received: %s\n", msg)
		}
	}()

	// Горутина, где отправляются сообщения на сервер.
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			// Чтение строки в байтах.
			s, err := reader.ReadBytes('\n')
			if err != nil {
				log.Printf("Error when reading from stdin: %v", err)
				continue
			}
			// Если пользователь ввел ctrl+c, то закрываем канал exitCh и завершаем работу горутины.
			if s[0] == 4 && len(s) == 3 {
				log.Println("Stop sending messages...")
				close(exitCh)
				return
			}
			// Посылаем введенное сообщение на сервер.
			if err := conn.WriteMessage(websocket.TextMessage, s[:len(s)-2]); err != nil {
				log.Printf("Error when sending message to server: %v", err)
			}
		}
	}()

	// Ожидание закрытия канала.
	<-exitCh
	log.Println("Exit program...")

	// output example:
	// 2022/03/10 19:24:03 Connecting to ws://ws.ifelse.io:80...
	// 2022/03/10 19:24:03 Connected to ws://ws.ifelse.io:80
	// 2022/03/10 19:24:03 Received: Request served by d7e94330
	// abcd
	// 2022/03/10 19:24:08 Received: abcd
	// абвг
	// 2022/03/10 19:24:14 Received: абвг
	// hello!
	// 2022/03/10 19:24:19 Received: hello!
	// ^D
	// 2022/03/10 19:24:22 Stop sending messages...
	// 2022/03/10 19:24:22 Exit program...
}
