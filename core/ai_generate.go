// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"bytes"         // Added import for bytes
	"encoding/json" // Added import for json
	"io/ioutil"     // Deprecated; replace with io and os if needed
	"log"           // Added import for log
	"net/http"
	"os"
	"strconv"

	"github.com/venusgalstar/go-ethereum/core/types"
)

func GenerateAI(tx *types.Transaction, msg *Message) {

	server := os.Getenv("AI_SERVER_IP")
	port := os.Getenv("AI_SERVER_PORT")

	if port == "" {
		port = "3000" // Default value
	}
	if server == "" {
		server = "127.0.0.1"
	}

	url := "http://" + server + ":" + port + "/generate"

	data := map[string]string{
		"hash":  tx.Hash().Hex(),
		"from":  msg.From.Hex(),
		"to":    msg.To.Hex(),
		"nonce": strconv.FormatUint(msg.Nonce, 10),
		"value": strconv.FormatUint(msg.Value.Uint64(), 10),
		"data":  string(msg.Data),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print the response
	log.Printf("Response Status: %s\n", resp.Status)
	log.Printf("Response Body: %s\n", string(body))
}

func GetGenerated(txHash string) (inscription string) {

	server := os.Getenv("AI_SERVER_IP")
	port := os.Getenv("AI_SERVER_PORT")

	if port == "" {
		port = "3000" // Default value
	}
	if server == "" {
		server = "127.0.0.1"
	}

	url := "http://" + server + ":" + port + "/getGenerated"

	data := map[string]string{
		"hash": txHash,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print the response
	log.Printf("Response Status: %s\n", resp.Status)
	return string(body)
}
