package main

import (
	"encoding/json"
	"fmt"
	"interpreter/evaluator"
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"net/http"
	"strings"
)

type CodePayload struct {
	Code string `json:"code"`
}

func interpreterEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var payload CodePayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		output := readInput(payload.Code)
		response := map[string]interface{}{
			"message": output,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

var output string

func readInput(input string) string {
	lines := strings.Split(input, ";")
	env := object.NewEnvironment()

	for _, line := range lines {
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(output, p.Errors())
		}
		evaluator.Eval(program, env)
	}
	return output
}

func printParserErrors(out string, errors []string) {
	output := ""
	output += "Parser Errors: \n"
	for _, msg := range errors {
		output += fmt.Sprintf("\t" + msg + "\n")
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	interpreterEndpoint(w, r)
	http.ListenAndServe(":3000", nil)
}
