package lib

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"path/filepath"
)

func LoadDict() dictInterface {
	start := time.Now()

	fmt.Println("loading dictionary...")

	dict := &Dictionary{definitions: make(map[string]*Definition)}

	dict.setFolder("data/old/")

	for ch := 'A'; ch <= 'Z'; ch++ {
		dict.loadData("wrangle/cleaned/" + string(ch) + ".json")
	}

	dict.PrintSize()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)

	fmt.Println()

	return dict
}

func LoadLLMDict() dictInterface {
	start := time.Now()

	fmt.Println("loading dictionary...")

	dict := &Dictionary{definitions: make(map[string]*Definition)}

	dict.setFolder("data/llmgen/")

	dict.loadData("wrangle/llmgen/gd.json")
	
	dict.PrintSize()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)

	fmt.Println()

	return dict
}

func LoadJSONDict(infile string, outfol string) dictInterface {
	start := time.Now()

	fmt.Println("loading dictionary...")

	dict := &Dictionary{definitions: make(map[string]*Definition)}

	dict.setFolder(outfol)

	dict.loadData(infile)
	
	dict.PrintSize()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)

	fmt.Println()

	return dict
}

func LoadWNDict() dictInterface {
	start := time.Now()

	fmt.Println("loading dictionary...")

	dict := &WNdict{definitions: make(map[string][]*WNdef), IDMappings: make(map[string]*WNdef)}

	dict.loadData("wn.json")

	dict.PrintSize()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)

	fmt.Println()

	return dict
}

func folderExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return fileInfo.IsDir()
}

func write(li []string, fn string) {
	json, err := json.MarshalIndent(li, "", " ")
	
	dir := filepath.Dir(fn)

	if (!folderExists(dir)){
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}	
	}
	

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		err = os.WriteFile(fn, json, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getNodes(fn string) []string {
	file, err := os.Open(fn)
	if err != nil {
		fmt.Println("error loading json")
		return []string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var txt string

	for scanner.Scan() {
		line := scanner.Text()
		txt = txt + line
	}

	bytes := []byte(txt)

	var myData []string

	json.Unmarshal(bytes, &myData)

	return myData
}

func Solve(dict dictInterface) {
	folder := dict.getFolder()

	tGraph := &Graph{vertices: make(map[string]*Vertex), pqMap: make(map[string]*Item)}

	dict.AddData(tGraph)
	tGraph.pqInit()

	listFree := tGraph.top()

	write(listFree, folder+"undefWords.json")

	start := time.Now()

	delNodes := tGraph.FVS()

	write(delNodes, folder+"delNodes.json")

	fmt.Println("nodes removed: ", len(delNodes))

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func reconstructWord(dict dictInterface, word string, fn string) {
	folder := dict.getFolder()

	delNodes := getNodes(folder + fn)

	defn := dict.getDef(word)

	fmt.Println(defn)

	defn = dict.expandDef(delNodes, word)

	fmt.Println(defn)
}

func exportSol(dict dictInterface, fn string, fn2 string) {
	start := time.Now()

	folder := dict.getFolder()

	delNodes := getNodes(folder + fn)

	m := dict.export(delNodes)

	b, err := json.MarshalIndent(m, "", "")

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		err = os.WriteFile("data/sol/"+fn2, b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func cullSolution(dict dictInterface, fn string) {
	folder := dict.getFolder()

	tGraph := &Graph{vertices: make(map[string]*Vertex)}

	dict.AddData(tGraph)

	listFree := tGraph.top()

	delNodes := getNodes(folder + fn)

	start := time.Now()

	cullNodes := tGraph.cullSol(delNodes, listFree)

	write(cullNodes, folder+"cullNodes.json")

	fmt.Println("nodes removed: ", len(cullNodes))

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func simulatedAnnealing(dict dictInterface, fn string) {
	folder := dict.getFolder()

	tGraph := &Graph{vertices: make(map[string]*Vertex)}

	dict.AddData(tGraph)

	listFree := tGraph.top()

	delNodes := getNodes(folder + fn)

	start := time.Now()

	simNodes := tGraph.simAnneal(delNodes, listFree)

	write(simNodes, folder+"simNodes.json")

	fmt.Println("nodes removed: ", len(simNodes))

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func GraphVerify(dict dictInterface, fn string) {
	folder := dict.getFolder()

	delNodes := getNodes(folder + fn)

	tGraph := &Graph{vertices: make(map[string]*Vertex)}

	dict.AddData(tGraph)

	listFree := tGraph.top()

	start := time.Now()

	verified := tGraph.verify(delNodes, listFree)

	fmt.Println("verified: ", verified)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func alternateVerify(dict dictInterface, fn string) {
	folder := dict.getFolder()

	delNodes := getNodes(folder + fn)

	tGraph := &Graph{vertices: make(map[string]*Vertex)}

	dict.AddData(tGraph)

	start := time.Now()

	tGraph.popAlgo()

	for _, k := range delNodes {
		delList := tGraph.DeleteVertex(k)

		pops, delList := tGraph.popList(delList)

		for pops != 0 {
			pops, delList = tGraph.popList(delList)
		}
	}

	fmt.Println(tGraph.Size())

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func dictVerify(dict dictInterface, fn string) {
	start := time.Now()

	folder := dict.getFolder()

	delNodes := getNodes(folder + fn)

	verified := dict.verify(delNodes)

	fmt.Println("verified: ", verified)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

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

func exportTrees(dict dictInterface, fn string) {
	folder := dict.getFolder()

	delNodes := getNodes(folder + fn)

	tGraph := &Graph{vertices: make(map[string]*Vertex)}
	dict.AddData(tGraph)

	fmt.Println("exporting trees...")

	var export map[string]expGraph = make(map[string]expGraph)

	for _, v := range tGraph.vertices {
		if strings.Contains(v.key, "/") {
			continue
		}

		var set []string
		set = append(set, v.key)

		var X []string

		var g expGraph

		for len(set) != 0 {
			key := set[0]
			set = append(set[:0], set[1:]...)

			var b bool = false
			for _, x := range X {
				if key == x {
					b = true
					break
				}
			}
			if b {
				continue
			} else {
				X = append(X, key)
			}

			vert := tGraph.vertices[key]
			n := node{vert.key}
			g.Nodes = append(g.Nodes, n)

			b = false
			for _, del := range delNodes {
				if vert.key == del {
					b = true
				}
			}
			if b {
				continue
			}

			for _, neighbor := range vert.inList {
				set = append(set, neighbor.key)

				l := link{vert.key, neighbor.key}

				g.Links = append(g.Links, l)
			}

		}

		export[v.key] = g

		b, err := json.MarshalIndent(export, "", " ")

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		} else {
			str := folder + "trees/" + v.key + ".json"
			err = os.WriteFile(str, b, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		export = make(map[string]expGraph)

	}

}

func exportNames(dict dictInterface) {
	folder := dict.getFolder()

	export := dict.getNames()

	b, err := json.MarshalIndent(export, "", "")

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		err = os.WriteFile(folder+"names.json", b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func exportJson(dict dictInterface) {
	folder := dict.getFolder()

	tGraph := &Graph{vertices: make(map[string]*Vertex)}
	dict.AddData(tGraph)

	fmt.Println("exporting graph...")

	var export expGraph

	for _, vert := range tGraph.vertices {
		n := node{vert.key}

		export.Nodes = append(export.Nodes, n)

		for _, out := range vert.outList {
			var l link

			if out.key != "" {
				l = link{vert.key, out.key}

				export.Links = append(export.Links, l)
			}

		}

	}

	b, err := json.MarshalIndent(export, "", "")

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		err = os.WriteFile(folder+"expJson.json", b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func exportCSV(dict dictInterface, fn string) {
	folder := dict.getFolder()

	rows := [][]string{
		{"source", "target"},
	}

	tGraph := &Graph{vertices: make(map[string]*Vertex)}
	dict.AddData(tGraph)

	if fn != "" {
		delNodes := getNodes(folder + fn)
		for _, k := range delNodes {
			tGraph.DeleteVertex(k)
		}
	}

	for _, vert := range tGraph.vertices {

		for _, out := range vert.outList {
			if out.key != "" {
				rows = append(rows, []string{vert.key, out.key})
			}
		}

	}

	csvfile, err := os.Create(folder + "expCSV.csv")

	if err != nil {
		log.Fatalf("Failed to create file, : %s", err)
	}

	cswriter := csv.NewWriter(csvfile)

	for _, row := range rows {
		_ = cswriter.Write(row)
	}

	cswriter.Flush()
	csvfile.Close()
}
