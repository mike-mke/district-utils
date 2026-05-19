package main

import (
	"district/utils"
	"fmt"
	"os"
	"strings"
)

const EXCEL_FILENAME = "./resources/sample-meetings.xlsx"

func main() {
	filename := getFilename()
	if len(os.Args) <= 1 {
		printUsage()
	} else {
		switch strings.ToLower(os.Args[1]) {
		case "-list":
			printList(filename)
		case "-javascript":
			printJavaScript(filename)
		case "-json":
			printJson(filename)
		case "-html":
			printHtml(filename)
		case "-email":
			printEmail(filename)
		default:
			printUsage()
		}
	}
}

func getFilename() string {
	if len(os.Args) >= 3 {
		return os.Args[len(os.Args)-1]
	}
	return EXCEL_FILENAME
}

func printList(filename string) {
	groupMap, keys := utils.ReadDistrictSpreadsheet(filename)
	if groupMap != nil && keys != nil {
		for _, key := range keys {
			group := groupMap[key]
			group.PrintItem()
			fmt.Println()
		}
		fmt.Println("Items:", len(groupMap))
	}
}

func printHtml(filename string) {
	groupMap, keys := utils.ReadDistrictSpreadsheet(filename)
	if groupMap != nil && keys != nil {
		printHtmlHeader(getDistrictNumber(groupMap, keys))
		for _, key := range keys {
			group := groupMap[key]
			group.PrintHtml()
			fmt.Println()
		}
		printHtmlFooter()
	}
}

func printJavaScript(filename string) {
	groupMap, keys := utils.ReadDistrictSpreadsheet(filename)
	if groupMap != nil && keys != nil {
		fmt.Println("var ALL_MEETINGS = [")
		for idx, key := range keys {
			if idx > 0 {
				fmt.Printf(",\n")
			}
			group := groupMap[key]
			group.PrintJavaScript()
		}
		fmt.Println("\n];")
	}
}

func printJson(filename string) {
	groupMap, keys := utils.ReadDistrictSpreadsheet(filename)
	if groupMap != nil && keys != nil {
		fmt.Println("{ \"meetings\": [")
		for idx, key := range keys {
			if idx > 0 {
				fmt.Printf(",\n")
			}
			group := groupMap[key]
			group.PrintJson()
		}
		fmt.Println("\n]")
		fmt.Println("}")
	}
}

func printHtmlHeader(districtNumber string) {
	fmt.Println("<html>")
	fmt.Println("<head>")
	fmt.Println("</head>")
	fmt.Println("<body>")
	fmt.Println("<style>")
	fmt.Println("\t.aaGroup {\n\t\tmargin-top: 0px;\n\t\tmargin-right: 25px;\n\t\tmargin-bottom: 15px;\n\t\tmargin-left: 25px;\n\t}")
	fmt.Println("</style>")
	fmt.Println("<h1>All Groups in District " + districtNumber + "</h1>")
	fmt.Println("<br/>")
}

func printHtmlFooter() {
	fmt.Println("</body>")
	fmt.Println("</html>")
}

func getDistrictNumber(groupMap map[string]utils.AaGroup, keys []string) string {
	key := keys[0]
	group := groupMap[key]
	return group.GetDistrictNumber()
}

func printEmail(filename string) {
	groupMap, keys := utils.ReadDistrictSpreadsheet(filename)
	if groupMap != nil && keys != nil {
		for _, key := range keys {
			group := groupMap[key]
			group.PrintEmail()
			fmt.Println()
		}
	}
}

func printUsage() {
	fmt.Println("usage: ./district |-list|-html|-javascript|-json| -file {filename}")
	fmt.Println("example: ./district -list")
	fmt.Println("example: ./district -html -file ./resources/sample-meetings.xlsx")
}
