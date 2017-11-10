package main

import (
	"fmt"
	"net/http"
	"os"
	"bufio"
	"regexp"
)
//main handler that takes users input and returns the chat bots response
func userInputHandler(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("value");
	response := "";
	previousInput:="";
	
	//static general resaponses
	//responses if user-input is empty
	if (userInput == ""){
		response = "you dont say very much, are you feeling ok?";
	}else if (userInput == previousInput){ //responses when the user repeats them selfs
		response = "you seem to be repeating yourself, are you feeling ok?";
	}else{
		response=getResponse(userInput);
		if(response == "false"){//responses if eliza dose not recognise the users phrase
			response = "I dont quite understand what you mean, could you clarify what you mean?"
		}
	}
	fmt.Fprintf(w, "%s", response)
	previousInput = userInput;
}

func getResponse(s string)(string){
	//input := s;
	
	type key struct {
    keyword string
    responses  []string
	}
	keyList:=[]key{}; 
	response :="false";
	index :=0;
	firstScan := true;
	file, _ := os.Open("eliza.dat");
	defer file.Close();
	scanner := bufio.NewScanner(file);
	newKey := key{keyword:"",responses:[]string{}};
	for scanner.Scan() {
        matched,_:=regexp.MatchString("@Key@.*",string(scanner.Text()))
		if(matched==true){
			index = 1;
			if(firstScan!=true){
				keyList = append(keyList,newKey)
			}else{
				firstScan=false;
			}
		}else{
			index++
		}
		if(index ==2){
			newKey = key{keyword:string(scanner.Text()),responses:[]string{}};
		}else if(index>2){
			newKey.responses = append(newKey.responses,string(scanner.Text()));
		}
    }
	
	for _, v := range keyList {
	fmt.Printf(v.keyword,v.responses);
	}
	
	return response;
}

func main() {

	
	fs := http.FileServer(http.Dir("Html"))
	http.Handle("/", fs)

	http.HandleFunc("/user-input", userInputHandler)
	http.ListenAndServe(":8080", nil)
}