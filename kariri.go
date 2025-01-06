package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
		panic("The platform is not supported. I am unable to clear the terminal screen.")
	}
}

func list_libs() {
	fmt.Println(string("\033[37m"), " Here is a list of the available libraries in the libs/ directory:", string("\033[0m"))
	entries, err := os.ReadDir("./libs/")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		fmt.Println("   -", string("\033[33m"), e.Name(), string("\033[0m"))
	}
}

// Struktura danych z JSON
type Docs struct {
	Lib  string `json:"lib"`
	Docs []Doc  `json:"docs"`
}

type Doc struct {
	Ph1 string `json:"ph1"`
	Ph2 string `json:"ph2"`
	Ph3 string `json:"ph3"`
}

func main() {

	js := os.Args[1:]

	if len(js) == 0 {
		fmt.Println(string("\033[37m"), "\n  It is required to enter the library name with *json extension", string("\033[0m"))
		fmt.Println("  e.g.:", string("\033[3m"), "./kariri libs/the_name_of_lib.json \n", string("\033[0m"))
		list_libs()
	} else {

		js := os.Args[1]

		jsonFile, err := os.Open(js)

		if err != nil {
			fmt.Println(err)
		}

		defer jsonFile.Close()

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		var docs Docs

		var (
			o int
			j []int
			s []string
		)

		// Losowanie, mieszanie :)
		r := rand.New(rand.NewSource(time.Now().Unix()))
		json.Unmarshal(byteValue, &docs)

		// START PROGRAMU: KARIRI
		for {

			j = r.Perm(len(docs.Docs))

			if len(j) == 0 {
				fmt.Println(string("\033[31m"), "  Probably the library or directory name was not entered correctly.", string("\033[0m"))
				os.Exit(0)
			} else {
				CallClear()
				fmt.Println(string("\033[33m"), ".::K:A:R:i:R:I::.\n", string("\033[0m"))
				fmt.Println(string("\033[32m"), "Library: ", string("\033[0m")+docs.Lib)
				fmt.Print(string("\033[32m"), " "+"Number of words: ", string("\033[0m"))
				fmt.Printf("%v \n\n", len(j))
			}

			fmt.Println(string("\033[31m"), "a:"+string("\033[37m"), "Remember", string("\033[0m"))
			fmt.Println(string("\033[31m"), "b:"+string("\033[37m"), "Check\n", string("\033[0m"))
			fmt.Println(string("\033[31m"), "q:"+string("\033[37m"), "Quit\n", string("\033[0m"))

			fmt.Print(">> ")
			reader := bufio.NewReader(os.Stdin)
			sw, err := reader.ReadByte()

			if err != nil {
				fmt.Println(err)
			}

			switch sw {

			case 'q':
				// CallClear()
				fmt.Println(string("\033[33m"), "\n.:: Thank you and wait to see you again :) ::.", string("\033[0m"))
				os.Exit(0)

			case 'b':
				CallClear()
				fmt.Println(string("\033[33m"), ".::K:A:R:i:R:I::.\n", string("\033[0m"))
				fmt.Println(string("\033[32m"), ".::Test your knowledge. See the words that cause problems. Repeat once more.::.\n", string("\033[0m"))

				j = r.Perm(len(docs.Docs))

				for i := 0; i < len(j); i++ {

					fmt.Print(string(" \033[37m"), docs.Docs[j[i]].Ph1+string("\033[36m"), " --> ")
					reader := bufio.NewReader(os.Stdin)
					odp, _ := reader.ReadString('\n')
					odp = strings.TrimSuffix(odp, "\n")
					if err != nil {
						fmt.Println(err)
					}
					ok := strings.Compare(odp, docs.Docs[j[i]].Ph2)

					if ok == 0 {
						o++
					} else {
						fmt.Println(string("\033[31m"), "should be: "+string("\033[36m"), docs.Docs[j[i]].Ph2, string("\033[0m"))
						s = append(s, docs.Docs[j[i]].Ph2)
					}
				}

				fmt.Printf("\n")
				fmt.Printf("Number of correct answers: %v / %v\n", o, len(j))
				fmt.Print("Repeat the words:\n")
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

				for i := 0; i < len(j); i++ {
					fmt.Println(string("\033[36m"), docs.Docs[j[i]].Ph1+string("\033[37m"), "  "+docs.Docs[j[i]].Ph2, "  "+string("\033[33m"), docs.Docs[j[i]].Ph3, string("\033[0m"))
					fmt.Scanln()
				}
				fmt.Println(string("\033[37m"), "\n These were all the examples in the indicated database.", string("\033[0m"))
				fmt.Scanln()
			default:
				fmt.Printf("\033[1A\033[K")
				fmt.Println(">> Incomprehensible command. Try again.")
				fmt.Scanln()
			}

		}
	}
}
