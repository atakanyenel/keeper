package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	hostFilename = ".kept"
	envFilename  = ".env"
	apiURL       = "http://keeperdatabase.herokuapp.com/hosts/"
)

func writeToFile(context, filename string) {
	err := ioutil.WriteFile(filename, []byte(context), 0644)
	if err != nil {
		panic(err)
	}
}
func readHostname() (string, error) {
	hostname, err := ioutil.ReadFile(hostFilename)

	return string(hostname), err
}

func getSingleHost(hostname string) string {
	resp, err := http.Get(apiURL + hostname)
	if err != nil {
		//handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	response := string(body)
	return response
}
func getAllHosts() string {
	resp, err := http.Get(apiURL)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	response := string(body)
	response = strings.TrimSpace(response)
	return response
}
func sendPost(hostname string) *http.Response {
	hc := http.Client{}
	form := url.Values{}
	form.Add("hostname", hostname)

	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(form.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := hc.Do(req)
	if err != nil {
		fmt.Printf("Error: %q", err)
		return resp
	}
	return resp

}
func printHelp() {
	fmt.Printf("Usage:\n---\n\nkeeper update\nkeeper tag <hostname>\nkeeper get <hostname>\n")
}
func main() {

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "tag":
		if len(os.Args) < 3 {
			printHelp()
			os.Exit(1)
		}
		fmt.Println("tagging")
		hostname := os.Args[2]
		oldHostname, err := readHostname()
		if err != nil {
			fmt.Printf("Creating new host: %s\n", hostname)
		} else {
			fmt.Printf("Changing %s to %s\n", oldHostname, hostname)
		}
		writeToFile(hostname, hostFilename)
		resp := sendPost(hostname)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
		fmt.Println(string(body))

	case "update":

		hostname, err := readHostname()
		if err == nil {
			sendPost(hostname)
		}
		hosts := getAllHosts()
		if len(hosts) < 1 {
			fmt.Println("No saved host")
			os.Exit(1)
		}
		fmt.Print(hosts)
		writeToFile(hosts, envFilename)
	case "get":
		if len(os.Args) < 3 {
			printHelp()
			os.Exit(1)
		}
		ipaddress := getSingleHost(os.Args[2])
		fmt.Printf(ipaddress)
	}
}
