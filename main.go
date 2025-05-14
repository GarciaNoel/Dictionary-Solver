package main

import (
	"github.com/garcianoel/dictionary-solver/lib"
)


func main() {

	//lib.handleServer("wnSol.json")

	// infile prefix here is ./wrangle/cleaned/
	dict := lib.LoadJSONDict("wrangle/llmgen/gd.json","data/llmgen/") 
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
