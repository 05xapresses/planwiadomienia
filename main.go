package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
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

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Wystąpił błąd. Możliwe, że jest podane więcej lekcji niż przerw - upewnij się, że konfiguracja jest prawidłowa.")
			fmt.Println()
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	cnt, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Nie znaleziono konfiguracji, lub nie udało jej się odczytać!")
		os.Exit(1)
		return
	}
	errJ := json.Unmarshal(cnt, &config)
	if errJ != nil {
		fmt.Println(errJ)
		return
	}
	if config.DlugoscLekcji < 0 {
		fmt.Println("Długość lekcji nie może być mniejsza niż 0!")
		os.Exit(1)
		return
	}
	pierwsza, err2 := time.Parse("03:04", config.PierwszaLekcja)
	if err2 != nil {
		fmt.Println("Nieprawidłowy format pierwszej lekcji!")
		fmt.Println(err2)
		os.Exit(1)
		return
	}
	for dzien, planLekcje := range config.Plan {
		rozpoczecie := pierwsza
		syf, blad := strconv.Atoi(dzien)
		if blad != nil || syf < 0 || syf > 6 {
			fmt.Printf("Nieprawidłowy numer dnia %v!\n", syf)
			os.Exit(1)
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

	fmt.Println("= Konfiguracja =")
	fmt.Print("długość lekcji: ")
	fmt.Print(config.DlugoscLekcji)
	fmt.Println(" minut")
	fmt.Print("godzina rozpoczęcia pierwszej lekcji: ")
	fmt.Println(pierwsza.Format("15:04"))
	fmt.Print("przerwy: ")
	fmt.Println(config.Przerwy)
	fmt.Println("wiadomość: ")
	fmt.Println(fmt.Sprintf(config.Wiadomosc, "*nazwa lekcji*", "*aktualna godzina*"))
	fmt.Print("webhook: ")
	fmt.Println(config.Webhook)
	fmt.Println("Odczytane lekcje:")
	keys := make([]int, 0, len(lekcje))
	for k := range lekcje {
		keys = append(keys, k)
	}
	sort.Ints(keys) // sortowanie dla normalnej kolejności
	for _, k := range keys {
		v := lekcje[k]
		fmt.Print("  ")
		fmt.Print(time.Weekday(k).String())
		fmt.Println(":")
		if len(v) == 0 {
			fmt.Println("    brak lekcji - można wyrzucić \"" + strconv.Itoa(k) + "\":[] z konfiguracji")
		}
		keys2 := make([]string, 0, len(v))
		for k := range v {
			keys2 = append(keys2, k)
		}
		sort.Strings(keys2) // sortowanie dla normalnej kolejności
		for _, godzina := range keys2 {
			lekcja := v[godzina]
			fmt.Print("    ")
			fmt.Print(godzina)
			fmt.Print(": ")
			fmt.Println(lekcja)
		}
	}
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
	_, err := http.Post(config.Webhook, "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Println("dzien dobry pragne przypomniec ze discord to syf i zawsze jest jakis problem")
		fmt.Println(err)
	}
}
