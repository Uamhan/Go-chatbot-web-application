package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//main handler that takes users input and returns the chat bots response
func userInputHandler(w http.ResponseWriter, r *http.Request) {
	//input from user taken from html page
	userInput := r.URL.Query().Get("value")
	//Seed the rand package with the current time.
	rand.Seed(time.Now().UnixNano())
	//Create a new instance of Eliza.
	eliza := ElizaFromFiles("data/responses.txt", "data/substitutions.txt")
	fmt.Fprintf(w, "%s", eliza.RespondTo(userInput))
}

// Replacer is a struct with two elements, a compiled regular expression and a string array containing replacements strings
type Replacer struct {
	original     *regexp.Regexp
	replacements []string
}

// ReadReplacersFromFile reads an array of Replacers from a text file.
// It takes a string which is the path to the data file.
func ReadReplacersFromFile(path string) []Replacer {
	// Open the file, logging a fatal error if it fails, close on return.
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Creates Replacers aray
	var replacers []Replacer

	// Read the file line by line.
	for scanner, readoriginal := bufio.NewScanner(file), false; scanner.Scan(); {
		// Read the next line and decide what to do.
		switch line := scanner.Text(); {
		// If the line starts with a # character then skip it.
		case strings.HasPrefix(line, "#"):
			// Do nothing
		// If we see a blank line, then make sure we indicate a new section.
		case len(line) == 0:
			readoriginal = false
		// If we haven't read the original, then append an element to the replacers array, compiling the regular expression.
		case readoriginal == false:
			replacers = append(replacers, Replacer{original: regexp.MustCompile(line)})
			readoriginal = true
		// Otherwise read a replacement and add it to the last replacer.
		default:
			replacers[len(replacers)-1].replacements = append(replacers[len(replacers)-1].replacements, line)
		}
	}
	// Return the replacers array.
	return replacers
}

// Eliza is a data structure representing a chatbot.
// The fields responses and substitutions are arrays of Replacers.
type Eliza struct {
	responses     []Replacer
	substitutions []Replacer
}

// ElizaFromFiles reads in text files containing responses and substituions
func ElizaFromFiles(responsePath string, substitutionPath string) Eliza {
	eliza := Eliza{}

	eliza.responses = ReadReplacersFromFile(responsePath)
	eliza.substitutions = ReadReplacersFromFile(substitutionPath)

	//returns instance of eliza with responses and substitutions loaded
	return eliza
}

// RespondTo takes a string as input and returns a string containg elizas response
func (me *Eliza) RespondTo(input string) string {
	// Look for a possible response.
	for _, response := range me.responses {
		// Check if the user input matches the original, capturing any groups.
		if matches := response.original.FindStringSubmatch(input); matches != nil {
			// Select a random response.
			output := response.replacements[rand.Intn(len(response.replacements))]
			// We'll tokenise the captured groups using the following regular expression.
			boundaries := regexp.MustCompile(`[\s,.?!]+`)
			// Fill the response with each captured group from the input.
			for m, match := range matches[1:] {
				// First split the captured group into tokens.
				tokens := boundaries.Split(match, -1)
				// Loop through the tokens.
				for t, token := range tokens {
					// Loop through the potential substitutions.
					for _, substitution := range me.substitutions {
						// Check if the original of the current substitution matches the token.
						if substitution.original.MatchString(token) {
							// If it matches, replace the token with one of the replacements Then break.
							tokens[t] = substitution.replacements[rand.Intn(len(substitution.replacements))]
							break
						}
					}
				}
				// Replace $1 with the first match, $2 with the second, etc.
				output = strings.Replace(output, "$"+strconv.Itoa(m+1), strings.Join(tokens, " "), -1)
			}
			// Send the filled answer back.
			return output
		}
	}
	// If there are no matches, then return this generic response.
	return "I dont quite understand, could you please clarify what you mean?"
}

func main() {
	//serves the html page to the user
	fs := http.FileServer(http.Dir("Html"))
	http.Handle("/", fs)

	//handles what happens when the user clicks on send
	http.HandleFunc("/user-input", userInputHandler)

	//what port the program listens and serves on.
	http.ListenAndServe(":8080", nil)
}
