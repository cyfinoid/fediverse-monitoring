package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"example.com/fediverse/Db"
)

var total_users int = 0
var total_posts int = 0
var total_comments int = 0
var url_nodes string = "https://nodes.fediverse.party/nodes.json"

func main() {
	startTime := time.Now()
	var wg sync.WaitGroup

	// ctx := context.TODO()
	// _, err := Db.ConnectToMongoDB(ctx)
	err := Db.Connect("dailygrowth")
	if err != nil {
		panic(err)
	}

	Db.CreateTableDG()
	tableName := Db.CreateTableNodes()
	errorTableName := Db.CreateTableError()

	if tableName == "" || errorTableName == "" {
		panic(err)
	} else {
		fmt.Println("Nodes Table name: " + tableName + "\nError Table name:" + errorTableName)
	}

	// url_node_info := "https://'$i'/.well-known/nodeinfo"

	body, err := GetApi(url_nodes)
	if err != nil {
		return
	}
	var data []interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("[NODES] Error unmarshalling jSON:", err)
		return
	}

	node_count := len(data)
	uniqueKeys := make(map[string]bool)

	numWorkers := 25 //no. of goroutines
	
	j := 0

	workerCh := make(chan struct{}, numWorkers)

	for _, node := range data {

		j++
		wg.Add(1)
		go process(tableName, errorTableName, node, uniqueKeys, j, &wg, workerCh, node_count)

	}

	wg.Wait()
	fmt.Println(node_count, total_users, total_comments, total_posts)
	elapsedTime := time.Since(startTime)
	fmt.Printf("[PSQL] Script execution time: %s\n", elapsedTime)
	Db.InsertDG(node_count, total_users, total_comments, total_posts, tableName, errorTableName, elapsedTime.String())

}

func process(tableName string, errorTableName string, node interface{}, uniqueKeys map[string]bool, j int, wg *sync.WaitGroup, workerCh chan struct{}, node_count int) {
	defer wg.Done()
	workerCh <- struct{}{}
	worker(tableName, errorTableName, node, uniqueKeys, j, node_count)
	<-workerCh
}

func worker(tableName string, errorTableName string, node interface{}, uniqueKeys map[string]bool, i int, node_count int) {

	body, err := GetApi("https://" + node.(string) + "/.well-known/nodeinfo")
	if err != nil {
		Db.InsertError(errorTableName, node.(string), "Couldn't fetch node info (TIMEOUT)")
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("[NODES] Error unmarshalling jSON:", err)
		return
	}

	tld, err := extractTLD(node.(string))
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	fmt.Println("[WORKER TLD]",tld)

	if linksInterface, ok := data["links"]; ok {

		if links, ok := linksInterface.([]interface{}); ok {
			href := links[0].(map[string]interface{})["href"].(string)

			if href == "" {
				return
			}
			var metadata map[string]interface{}
			body, err = GetApi(href)
			if err != nil {
				Db.InsertError(errorTableName, node.(string), "Couldn't fetch metadata (TIMEOUT)")
				return
			}
			err = json.Unmarshal(body, &metadata)
			if err != nil {
				fmt.Println("[METADATA] Error unmarshalling jSON:", err)
				Db.InsertError(errorTableName, node.(string), "Error unmarshalling metadata")

				return
			}

			fmt.Println("Iteration", i)

			// if data, ok := metadata["metadata"].(map[string]interface{}); ok {
			// 	fmt.Println("metadata present")
			// 	for key, value := range data {
			// 		if _, exists := uniqueKeys[key]; !exists {

			// 			uniqueKeys[key] = true
			// 			fmt.Printf("New Key: %s\n", key)
			// 			// writer.Write([]string{key})
			// 		}
			// 		if nestedMap, ok := value.(map[string]interface{}); ok {
			// 			processKeys(nestedMap, uniqueKeys)
			// 		}
			// 	}
			// } else {
			// 	fmt.Println("Empty metadata")
			// }

			softwareInterface := metadata["software"]
			var software map[string]interface{}

			if softwareInterface != nil {
				software, ok = softwareInterface.(map[string]interface{})
				if !ok {
					fmt.Println("[METADATA SOFTWARE] Error: Unable to convert 'software' value to map[string]interface{}.")
					Db.InsertError(errorTableName, node.(string), "Error converting metadata software value to map[string]interface{}")
				}
			} else {
				fmt.Println("Error: 'software' value is nil.")
				Db.InsertError(errorTableName, node.(string), "Metadata 'software' value is nil.")
				return
			}
			name := fmt.Sprintf("%s", software["name"])
			version := fmt.Sprintf("%s", software["version"])
			protocols := sortAndFilterProtocols(metadata["protocols"])
			ver, ok := metadata["version"].(string)
			if !ok {
				fmt.Println("Error: Unable to convert 'version' value to string.")
				Db.InsertError(errorTableName, name, "Error: Unable to convert 'version' value to string.")
			}
			regOpen := metadata["openRegistrations"]
			users, err := getFloat64(metadata, "usage", "users", "total")
			if err != nil {
				fmt.Println("[USERS] Error:", err)
				Db.InsertError(errorTableName, name, "Error calculating users float value")
				return
			}

			posts, err := getFloat64(metadata, "usage", "localPosts")
			if err != nil {
				fmt.Println("[POSTS] Error:", err)
				Db.InsertError(errorTableName, name, "Error calculating posts float value")
				return
			}

			totalUsers := int(users)
			totalPosts := int(posts)
			totalComments, ok := func() (int, bool) {
				if usage, ok := metadata["usage"].(map[string]interface{}); ok {
					if localCommentsValue, ok := usage["localComments"]; ok {
						if comments, ok := localCommentsValue.(float64); ok {
							return int(comments), true
						}
					}
				}

				return 0, false
			}()

			total_users += totalUsers
			total_posts += totalPosts
			total_comments += totalComments

			finalEcho := fmt.Sprintf("name: %s\nprotocols: %s\nver: %s\nregOpen: %t \ntotalPosts: %d\ntotalComments: %d\ntotalUser: %d", name, protocols, ver, regOpen, totalPosts, totalComments, totalUsers)
			fmt.Println(finalEcho)
			Db.InsertNodes(tableName, name, version, tld, totalUsers, totalComments, totalPosts, metadata)
		}

	} else {
		return
	}

}

//
// date&time+domain as key
//
// Try NoSQL, key= data-server, value = json, compare time
//

func GetApi(url string) ([]byte, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	if url == url_nodes {
		ctx, cancel = context.WithTimeout(context.Background(), 25*time.Second)
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	}
	defer cancel()

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("[GET] Error creating request:", err)
		// errorData := map[string]interface{}{
		// 	"error": "Error creating request",
		// }
		return nil, err
	}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[GET] Error making request:", err)
		// errorData := map[string]interface{}{
		// 	"error": "timeout",
		// }
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	return body, nil
}

func sortAndFilterProtocols(protocols interface{}) string {
	var filteredProtocols []string

	// Check if protocols is a slice
	if protocolsSlice, ok := protocols.([]interface{}); ok {
		// Iterate over the slice
		for _, p := range protocolsSlice {
			// Check if the element is a string
			if protocol, ok := p.(string); ok {
				// Filter protocols
				if protocol != "inbound" && protocol != "outbound" {
					filteredProtocols = append(filteredProtocols, protocol)
				}
			} else {
				// Handle unexpected type in the slice if needed
				fmt.Printf("Unexpected type in slice: %T\n", p)
			}
		}

		// Sort and format protocols
		sort.Strings(filteredProtocols)
		protocolsString := strings.Join(filteredProtocols, ": ")
		return protocolsString
	}

	// Handle unexpected type for protocols
	fmt.Printf("Unexpected type: %T\n", protocols)
	return ""
}

func extractTLD(input string) (string, error) {
	// Define a regular expression for matching TLDs
	tldRegex := regexp.MustCompile(`[a-zA-Z0-9-]+\.([a-zA-Z]{2,})$`)

	// Find the first match
	matches := tldRegex.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", fmt.Errorf("No TLD found in the input string")
	}

	return matches[1], nil
}

func getFloat64(m map[string]interface{}, keys ...string) (float64, error) {
	var value interface{} = m
	for _, key := range keys {
		nestedMap, ok := value.(map[string]interface{})
		if !ok {
			return 0, fmt.Errorf("Key %s not found or has an invalid type", key)
		}
		value = nestedMap[key]
		if value == nil {
			return 0, fmt.Errorf("Key %s has a nil value", key)
		}
	}

	floatValue, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("Unable to convert value to float64")
	}
	return floatValue, nil
}
