package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	jira "github.com/andygrunwald/go-jira"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input-file>")
		os.Exit(1)
	}

	apiToken := os.Args[1]

	// Set the Jira API endpoint and issue key
	apiEndpoint := "https://yocreoab.atlassian.net"

	user := "per.e.olofsson@gmail.com"

	tp := jira.BasicAuthTransport{
		Username: user,
		Password: apiToken,
	}

	client, _ := jira.NewClient(tp.Client(), apiEndpoint)

	input := "test.adoc"

	output := "testout.adoc"

	content, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}

	// Define the regular expression pattern
	pattern := `<CC-(\d+)[^>]*>`

	// Create a regular expression object
	r := regexp.MustCompile(pattern)

	// Find all instances of the pattern in the content
	matches := r.FindAllStringSubmatchIndex(string(content), -1)

	// Loop through the matches and replace each one
	newContent := ""
	lastIndex := 0
	for _, match := range matches {
		startIndex := match[0]
		endIndex := match[1]

		// Append the text before the match
		newContent += string(content[lastIndex:startIndex])

		// Extract the NUM value from the match
		numString := string(content[match[2]:match[3]])

		issue, _, err := client.Issue.Get("CC-"+numString, nil)
		if err != nil {
			panic(err)
		}

		newContent += "<" + issue.Key + ":" + issue.Fields.Summary + ">"

		lastIndex = endIndex
	}

	// Append the remaining text after the last match
	newContent += string(content[lastIndex:])

	// Print the updated content
	fmt.Println(newContent)

	// Write the updated content to the output text file
	err = ioutil.WriteFile(output, []byte(newContent), 0644)
	if err != nil {
		panic(err)
	}

}
