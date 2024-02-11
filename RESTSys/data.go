package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
	"strconv"
)

type DataPayload struct {
	CPU    int `json:"cpu"`
	Memory int `json:"memory"`
	Disk   int `json:"disk"`
}
type Statistic struct {
	Name string `json:"name"`
	Mean int    `json:"mean"`
	Std  int    `json:"std"`
}

type DataPayloadConfig []Statistic

func main() {
	// Générer des nombres aléatoires pour les données
	var configObject DataPayloadConfig
	args := os.Args

	config, err := os.ReadFile("config.txt")
	if err != nil {
		config = []byte("[{\"name\":\"cpu\",\"mean\":3,\"std\":1},{\"name\":\"memory\",\"mean\":30,\"std\":10}]")
	}
	fmt.Println("config: " + string(config) + "\n")
	json.Unmarshal(config, &configObject)

	rand.Seed(time.Now().UnixNano())
	cpu := rand.Intn(11)
	memory := rand.Intn(101)
	disk := rand.Intn(1001)

	for i := 0; i < len(configObject); i++ {
		if configObject[i].Name == "cpu" {
			cpu = int(math.Round(math.Abs(rand.NormFloat64()*float64(configObject[i].Std) + float64(configObject[i].Mean))))
		} else if configObject[i].Name == "memory" {
			memory = int(math.Round(math.Abs(rand.NormFloat64()*float64(configObject[i].Std) + float64(configObject[i].Mean))))
		}
	}
	for i:=1; i<len(args); i++ {
		switch (i) {
		case 1:
			cpu, _ = strconv.Atoi(args[i])
		case 2:
			memory,_ = strconv.Atoi(args[i])
		case 3:
			disk,_ = strconv.Atoi(args[i])
		}

	}
	payload := DataPayload{
		CPU:    cpu,
		Memory: memory,
		Disk:   disk,
	}
	fmt.Printf("%+v\n", payload)
	// Convertir la charge utile en JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Erreur lors de la conversion de la charge utile en JSON : %v", err)
	}

	// Envoyer une requête POST à l'API REST
	resp, err := http.Post("http://closedloop-v2-monitoring-deployment-service.com:80/send-data", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Erreur lors de l'envoi de la requête POST : %v", err)
	}
	defer resp.Body.Close()

	// Vérifier la réponse
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("La requête a échoué avec le code de statut : %d", resp.StatusCode)
	}

	// Lecture de la réponse
	var response struct {
		Message string `json:"message"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture de la réponse : %v", err)
	}

	// Afficher le message de réponse
	fmt.Println("Réponse de l'API :", response.Message)
}
