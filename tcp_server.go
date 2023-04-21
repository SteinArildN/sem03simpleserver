package main

import (
	"io"
	"log"
	"net"
	"sync"
	"strings"
	"fmt"
	"strconv"

	"github.com/SteinArildN/funtemps/conv"
	"github.com/SteinArildN/is105sem03/mycrypt"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:16")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03) - 4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
					//log.Println(buf[:n])
					switch msg := string(dekryptertMelding[:n]); msg {
					//log.Println(dekryptertMelding)
					//switch msg := string(buf[:n]); msg {
  				        case "ping":
						//_, err = c.Write([]byte("pong"))
						kryptertPong := mycrypt.Krypter([]rune(string("pong")), mycrypt.ALF_SEM03, 4)
						_, err = c.Write([]byte(string(kryptertPong[:n])))
					case "Kjevik;SN39040;18.03.2022 01:50;6":
						////fahr := 0
						elementArray := strings.Split(msg, ";")
						celsius, _ := strconv.ParseFloat(elementArray[3], 64)
						fahr := conv.CelsiusToFahrenheit(celsius)
						newEArray := []string{elementArray[0], elementArray[1], elementArray[2], fmt.Sprintf("%v",fahr)}
						log.Println(newEArray)
						//fahr := conv.CelsiusToFahrenheit(celsius)
						finalString := strings.Join(newEArray, ";")// + ";" + fmt.Sprintf("%v", fahr)
						log.Println(finalString)
						//_, err c.Write(elementArray[0] + elementArray[1] + elementArray[2] + fahr)
						kryptertString := mycrypt.Krypter([]rune(string(finalString)), mycrypt.ALF_SEM03, 4)
						log.Println(string(kryptertString))
						_, err = c.Write([]byte(string(kryptertString[:len(kryptertString)])))
						//c.Write([]byte(finalString))
						//_, err = c.Write([]byte("jippi"))
					default:
						_, err = c.Write(buf[:n])
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
