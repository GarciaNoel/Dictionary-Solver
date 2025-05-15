package main

import (
	"flag"
	"github.com/garcianoel/dictionary-solver/lib"
)

var (
    infile    string
    outfol    string
)

func init() {
    flag.StringVar(&infile, "jsonpath", "wrangle/llmgen/gd.json", "path to json dictionary")
    flag.StringVar(&outfol, "folderpath", "data/llmgen/", "path to directory for data outuput")
}


func main() {

	flag.Parse()

	//lib.handleServer("wnSol.json")

	dict := lib.LoadJSONDict(infile, outfol) 
	//dict := lib.LoadLLMDict()
	//dict := lib.LoadWNDict()

	lib.Solve(dict)

	//lib.reconstructWord(dict, "happy", "delNodes.json")

	//lib.exportSol(dict, "delNodes.json", "oldSol.json")

	//lib.simulatedAnnealing(dict, "delNodes.json")

	//lib.cullSolution(dict, "delNodes.json")

	lib.GraphVerify(dict, "delNodes.json")

	//lib.alternateVerify(dict, "delNodes.json")

	//dictVerify(dict, "cullNodes.json")

	//exportTrees(dict, "delNodes.json")

	//lib.exportNames(dict)

	//lib.exportJson(dict)

	//lib.exportCSV(dict, "delNodes.json")
}
