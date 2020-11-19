package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type cfg struct {
	DlugoscLekcji  int
	PierwszaLekcja string
	Przerwy        []int
	Plan           map[string][]string
	Wiadomosc      string
	Webhook        string
}

var config = cfg{}
var ostatnioWyslano string
var lekcje = make(map[int]map[string]string)

func main() {

	cnt, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Nie znaleziono konfiguracji!")
		return
	}
	errJ := json.Unmarshal(cnt, &config)
	if errJ != nil {
		fmt.Println(errJ)
		return
	}
	if config.DlugoscLekcji < 0 {
		log.Fatalln("Długość lekcji nie może być mniejsza niż 0!")
		return
	}
	fmt.Println(config)
	pierwsza, err2 := time.Parse("03:04", config.PierwszaLekcja)
	fmt.Println(pierwsza)
	if err2 != nil {
		fmt.Println(err2)
		log.Fatalln("Nieprawidłowy format pierwszej lekcji!")
		return
	}
	for dzien, planLekcje := range config.Plan {
		rozpoczecie := pierwsza
		syf, blad := strconv.Atoi(dzien)
		if blad != nil || syf < 1 || syf > 5 {
			log.Fatalln("Nieprawidłowy numer dnia!")
			return
		}
		lekcje[syf] = make(map[string]string)
		i := 0
		for _, lekcja := range planLekcje {
			if len(lekcja) != 0 {
				lekcje[syf][rozpoczecie.Format("15:04")] = lekcja
			}
			rozpoczecie = rozpoczecie.Add(time.Duration(config.DlugoscLekcji) * time.Minute)
			rozpoczecie = rozpoczecie.Add(time.Duration(config.Przerwy[i]) * time.Minute)
			i++
		}
	}
	fmt.Println(lekcje)

	for {
		dzien := time.Now().Weekday()
		godzina := time.Now().Format("15:04")
		dzienObj, dzienIstnieje := lekcje[int(dzien)]
		if dzienIstnieje {
			lekcja, godzinaIstnieje := dzienObj[godzina]
			if godzinaIstnieje {
				terazWyslano := strconv.Itoa(int(dzien)) + godzina
				if ostatnioWyslano != terazWyslano {
					ostatnioWyslano = terazWyslano
					sendmessage(lekcja, godzina)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}

}

func sendmessage(lekcja string, godzina string) {
	body, _ := json.Marshal(map[string]interface{}{
		"content": fmt.Sprintf(config.Wiadomosc, lekcja, godzina),
	})
	_, error := http.Post(config.Webhook, "application/json", bytes.NewReader(body))
	if error != nil {
		fmt.Println("dzien dobry pragne przypomniec ze discord to syf i zawsze jest jakis problem")
		fmt.Println(error)
	}
}
