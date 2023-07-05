package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Platforma nie jest wspierana. Nie mogę wyczyścić ekranu terminala.")
	}
}

func lista() {
	fmt.Println(string("\033[37m"), "  Oto lista dostępnych bibliotek w katalogu /json:", string("\033[0m"))
	entries, err := os.ReadDir("./json/")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		fmt.Println("   -", string("\033[33m"), e.Name(), string("\033[0m"))
	}
}

// Struktura danych z JSON
type Users struct {
	Baza  string `json:"base"`
	Users []User `json:"words"`
}

type User struct {
	PL  string `json:"pl"`
	ENG string `json:"eng"`
	ZDN string `json:"zdn"`
}

func main() {

	js := os.Args[1:]

	if len(js) == 0 {
		fmt.Println(string("\033[37m"), "\n   Wymagane jest wprowadzenie nazwy bazy danych z rozszerzeniem *json wraz z lokacją,", string("\033[0m"))
		fmt.Println("   np.:", string("\033[3m"), "./kariri json/<nazwa_bazy>.json \n", string("\033[0m"))
		lista()
	} else {

		js := os.Args[1]

		jsonFile, err := os.Open(js)

		if err != nil {
			fmt.Println(err)
		}

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var users Users

		var (
			o int
			j []int
			s []string
		)

		// Losowanie, mieszanie :)
		r := rand.New(rand.NewSource(time.Now().Unix()))
		json.Unmarshal(byteValue, &users)

		// START PROGRAMU: KARIRI
		for {

			j = r.Perm(len(users.Users))

			if len(j) == 0 {
				fmt.Println(string("\033[31m"), "\n   Prawdopodobnie nazwa bazy danych lub katalogu nie została wprowadzona poprawnie.", string("\033[0m"))
				os.Exit(0)
			} else {
				CallClear()
				fmt.Println(string("\033[33m"), ".::K:A:R:i:R:I::.\n", string("\033[0m"))
				fmt.Println(string("\033[32m"), "Biblioteka: ", string("\033[0m")+users.Baza)
				fmt.Print(string("\033[32m"), " "+"Liczba słów: ", string("\033[0m"))
				fmt.Printf("%v \n\n", len(j))
			}

			fmt.Println(string("\033[31m"), "a:"+string("\033[37m"), "Zapamiętaj", string("\033[0m"))
			fmt.Println(string("\033[31m"), "b:"+string("\033[37m"), "Sprawdź się\n", string("\033[0m"))
			fmt.Println(string("\033[31m"), "q:"+string("\033[37m"), "Wyjście\n", string("\033[0m"))

			fmt.Print(">> ")
			reader := bufio.NewReader(os.Stdin)
			sw, err := reader.ReadByte()

			if err != nil {
				fmt.Println(err)
			}

			switch sw {

			case 'q':
				// CallClear()
				fmt.Println(string("\033[33m"), "\n.:: Dziękuję i zapraszam ponownie :) ::.", string("\033[0m"))
				os.Exit(0)

			case 'b':
				CallClear()
				fmt.Println(string("\033[33m"), ".::K:A:R:i:R:I::.\n", string("\033[0m"))
				fmt.Println(string("\033[32m"), ".::Sprawdź swoją wiedzę. Zobacz słowa, które sprawiają problemy. Powtórz raz jeszcze.::.\n", string("\033[0m"))

				j = r.Perm(len(users.Users))

				for i := 0; i < len(j); i++ {

					fmt.Print(string(" \033[37m"), users.Users[j[i]].PL+string("\033[36m"), " --> ")
					reader := bufio.NewReader(os.Stdin)
					odp, _ := reader.ReadString('\n')
					odp = strings.TrimSuffix(odp, "\n")
					if err != nil {
						fmt.Println(err)
					}
					ok := strings.Compare(odp, users.Users[j[i]].ENG)

					if ok == 0 {
						o++
					} else {
						fmt.Println(string("\033[31m"), "powinno być: "+string("\033[36m"), users.Users[j[i]].ENG, string("\033[0m"))
						s = append(s, users.Users[j[i]].ENG)
					}
				}

				fmt.Printf("\n")
				fmt.Printf("Liczba prawidłowych odpowiedzi: %v / %v\n", o, len(j))
				fmt.Print("Powtórz słowa:\n")
				for i := 0; i < len(s); i++ {
					fmt.Printf("- %v\n", s[i])
				}

				s = nil
				j = nil
				o = 0
				fmt.Scanln()

			case 'a':
				CallClear()
				fmt.Println(string("\033[33m"), ".::K:A:R:i:R:I::.\n", string("\033[0m"))
				fmt.Println(string("\033[32m"), ".::Rozluźnij się. Przeczytaj wiele razy słowa i zdania. Zapamiętaj!::.\n", string("\033[0m"))

				for i := 0; i < len(j); i++ {
					fmt.Println(string("\033[36m"), users.Users[j[i]].ENG+string("\033[37m"), "  "+users.Users[j[i]].PL, "  "+string("\033[33m"), users.Users[j[i]].ZDN, string("\033[0m"))
					fmt.Scanln()
				}
				fmt.Println(string("\033[37m"), "\n Były to wszystkie przykłady we wskazanej bazie.", string("\033[0m"))
				fmt.Scanln()
			default:
				fmt.Printf("\033[1A\033[K")
				fmt.Println(">> Niezrozumiałe polecenie. Spróbuj raz jeszcze.")
				fmt.Scanln()
			}

		}
	}
}

// Zastosowanie keyloggera w opcji menu
