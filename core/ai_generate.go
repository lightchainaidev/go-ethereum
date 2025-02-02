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
   	"fmt"
    "io"
	"net/http"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

func GenerateAI(tx *types.Transaction) (string, error) {
	// Validate input
	if tx == nil {
		return "", fmt.Errorf("transaction cannot be nil")
	}

	// Fetch environment variables
	AI_SERVER_URL := os.Getenv("AI_SERVER_URL")

	if AI_SERVER_URL == "" {
		AI_SERVER_URL = "http://127.0.0.1:3000"
	}

	url:= AI_SERVER_URL

	to := ""
	if tx.To() != nil {
		to = tx.To().Hex()
	}

	// Prepare the request payload
	data := map[string]string{
		"hash":  tx.Hash().Hex(),
		"from":  params.SystemAddress.Hex(),
		"to":    to,
		"nonce": strconv.FormatUint(tx.Nonce(), 10),
		"value": strconv.FormatUint(tx.Value().Uint64(), 10),
		"data":  string(tx.Data()),
	}

	// Marshal the payload to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("AI error marshaling JSON: %w", err)
	}

	// Create a context with timeout to avoid hanging indefinitely
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// Create an HTTP request with the context
	// req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("AI error creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("AI error making POST request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("AI error reading response body: %w", err)
	}

	// Parse the JSON response
	var result struct {
		Data []struct {
			Text string `json:"text"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("AI error unmarshaling JSON response: %w", err)
	}

	// Check if the "data" array is not empty
	if len(result.Data) == 0 {
		return "", fmt.Errorf("AI no data found in response")
	}

	// Return the "text" field from the first item in the "data" array
	return result.Data[0].Text, nil
}

