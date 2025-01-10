package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	cmd := flag.String("cmd", "", "Command to execute: list, get, add, delete, health")
	id := flag.Int("id", 0, "ID of the item (for get or delete)")
	name := flag.String("name", "", "Name of the item (for add)")
	endpoint := flag.String("endpoint", "http://localhost:8080", "Server endpoint")
	flag.Parse()

	switch *cmd {
	case "list":
		listItems(*endpoint)
	case "get":
		if *id == 0 {
			fmt.Println("Please provide an ID for get command")
			os.Exit(1)
		}
		getItem(*endpoint, *id)
	case "add":
		if *name == "" {
			fmt.Println("Please provide a name for add command")
			os.Exit(1)
		}
		addItem(*endpoint, *name)
	case "delete":
		if *id == 0 {
			fmt.Println("Please provide an ID for delete command")
			os.Exit(1)
		}
		deleteItem(*endpoint, *id)
	case "health":
		healthCheck(*endpoint)
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

func listItems(endpoint string) {
	resp, err := http.Get(fmt.Sprintf("%s/api/items", endpoint))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func getItem(endpoint string, id int) {
	url := fmt.Sprintf("%s/api/items/%d", endpoint, id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func addItem(endpoint string, name string) {
	item := Item{Name: name}
	data, _ := json.Marshal(item)
	resp, err := http.Post(fmt.Sprintf("%s/api/items", endpoint), "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func deleteItem(endpoint string, id int) {
	url := fmt.Sprintf("%s/api/items/%d", endpoint, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Item deleted successfully")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
}

func healthCheck(endpoint string) {
	resp, err := http.Get(fmt.Sprintf("%s/api/health", endpoint))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
