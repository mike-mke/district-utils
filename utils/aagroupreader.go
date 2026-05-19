package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ReadDistrictSpreadsheet(filename string) (map[string]AaGroup, []string) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Get all the rows in the first sheet.
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	var headers []string
	groupMap := map[string]AaGroup{}
	keys := []string{}
	for idx, row := range rows {
		if idx == 0 {
			// populate headers
			for _, col := range row {
				headers = append(headers, strings.TrimSpace(col))
			}
		}

		if idx > 0 {
			group := CreateAaGroup(headers, row)
			key := group.GetKey()
			keys = append(keys, key)
			groupMap[key] = group
		}
	}

	sort.Strings(keys)
	return groupMap, keys
}
