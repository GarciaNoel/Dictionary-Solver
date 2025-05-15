package net

import (
	"github.com/gorilla/mux"
	"net/http"
	"io"
	"fmt"
	"os"
	"encoding/json"
	"log"
)

var SOL map[string][]string

type node struct {
	Name string `json:"name"`
}

type link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type expGraph struct {
	Nodes []node `json:"nodes"`
	Links []link `json:"links"`
}

func origHandler(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")

	val, ok := SOL[word]
	if ok {
		w.Write([]byte(val[0]))
	} else {
		w.Write([]byte(""))
	}

}

func newHandler(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")

	val, ok := SOL[word]
	if ok {
		w.Write([]byte(val[1]))
	} else {
		w.Write([]byte(""))
	}
}

func gHandler(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")

	file, err := os.Open("data/wn/trees/" + word + ".json")
	if err != nil {
		fmt.Print(err)
		return
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	var ret expGraph
	for {
		n, err := dec.Token()
		if err != nil {
			log.Fatal(err)
		}

		if n == word {
			var g expGraph
			if err := dec.Decode(&g); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			ret = g
			break
		} else {
			t, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			ctr := 1

			// slow implementation makes server spike 100% on single request!

			for ctr != 0 {
				t, err = dec.Token()
				if err != nil {
					log.Fatal(err)
				}

				var del json.Delim = '{'
				var open json.Token = del
				del = '}'
				var close json.Token = del

				if t == open {
					ctr += 1
				} else if t == close {
					ctr -= 1
				}
			}
		}
	}

	b, err := json.MarshalIndent(ret, "", " ")
	if err != nil {
		w.Write([]byte(""))
	}

	w.Write(b)

}

func HandleServer(fn string) {
	fmt.Println("starting server...")

	bytes, err := os.ReadFile("data/sol/" + fn)
	if err != nil {
		fmt.Print(err)
		return
	}

	json.Unmarshal(bytes, &SOL)

	bytes = nil

	r := mux.NewRouter()

	r.HandleFunc("/orig", origHandler).Methods("GET")
	r.HandleFunc("/new", newHandler).Methods("GET")
	r.HandleFunc("/graph", gHandler).Methods("GET")

	http.Handle("/", r)

	fmt.Println("server ready!")

	log.Fatal(http.ListenAndServe(":3001", nil))
}