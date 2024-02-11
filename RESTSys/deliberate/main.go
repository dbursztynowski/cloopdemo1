package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, r *http.Request) {
	type Statistic struct {
		Name string `json:"name"`
		Mean int    `json:"mean"`
		Std  int    `json:"std"`
	}

	type DataPayload []Statistic

	var dataPayload DataPayload
	var configObject DataPayload
	fmt.Printf("Hello\n")
	fmt.Printf(r.URL.Path + "\n")
	if r.URL.Path != "/hello/" {
		fmt.Printf("Error in Hello\n")
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Printf("GET\n")
		fmt.Fprintf(w, "HELLO\n")

	case "POST":
		config, err := os.ReadFile("../project/config.txt")
		if err != nil {
			//log.Fatal(err)
			config = []byte("[{\"name\":\"cpu\",\"mean\":3,\"std\":1},{\"name\":\"memory\",\"mean\":30,\"std\":10}]")
		}
		fmt.Println("config: " + string(config) + "\n")
		json.Unmarshal(config, &configObject)

		fmt.Printf("POST\n")
		reqBody, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintf(w, "%+v", string(reqBody))
		fmt.Printf("%+v\n", string(reqBody))
		json.Unmarshal(reqBody, &dataPayload)
		/*
		                err := json.NewDecoder(r.Body).Decode(&dataPayload)
		                if err != nil {
		                     fmt.Printf(err.Error(), http.StatusBadRequest)
				}*/
		/*for  _, configElement := range configObject {
			for  _, dataElement := range dataPayload {
				if dataElement.Name == configElement.Name {
					configElement.Mean = dataElement.Mean
					configElement.Std = dataElement.Std
				}
			}
		}*/
		for i := 0; i < len(configObject); i++ {
			for j := 0; j < len(dataPayload); j++ {
				if configObject[i].Name == dataPayload[j].Name {
					configObject[i].Mean = dataPayload[j].Mean
					configObject[i].Std = dataPayload[j].Std
				}
			}
		}
		f, err := os.Create("config.txt")
		if err != nil {
			log.Fatal(err)
		}
		// remember to close the file
		defer f.Close()

		// create new buffer
		buffer := bufio.NewWriter(f)
		configJson, _ := json.Marshal(configObject)
		configString := string(configJson[:])
		fmt.Printf("%+v\n", configString)
		_, err = buffer.WriteString(string(configString))
		//_, err = buffer.WriteString(string(reqBody))
		if err != nil {
			log.Fatal(err)
		}
		// flush buffered data to the file
		if err := buffer.Flush(); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/", hello)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
