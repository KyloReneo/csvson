package csvson

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func CSVToJSON(path string) {
	data, err := CSV2Slice(path)
	if err != nil {
		panic(fmt.Sprintf("error while handling csv file: %s\n", err))
	}
	res, _ := JSONString(data)
	WriteToJSON(res)
}

func JSONString(rows [][]string) (string, error) {

	// Split headers and values
	headerAttributes := rows[0]
	values := rows[1:]

	// Map of attributes and values
	entities := CSVRowsToMap(headerAttributes, values)

	// String of json formatted file
	json, err := MapToString(entities)
	if err != nil {
		log.Fatalf("Converting map of entities to json failed due to the error: %s", err)
	}
	return json, err
}

// Recives path of CSV file as input and returns a sliced array of rows
func CSV2Slice(CSVPath string) ([][]string, error) {
	inputFile := OpenCSV(CSVPath)
	slicedRows := ReadCSV(inputFile)
	return slicedRows, nil
}

// Opens the given path and returns the CSV file
func OpenCSV(path string) (CSVFile *os.File) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error while opening the path: %s\n %s", path, err)
		return nil
	}
	return file
}

// Reads the CSVFile and returns A 2D string slice of rows
func ReadCSV(CSVFile *os.File) (CSVRows [][]string) {
	var csvRows [][]string
	reader := csv.NewReader(CSVFile)
	for {
		// Read rows step by step
		row, err := reader.Read()
		// Stop reading when reached end of the file
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Somthing went wrong while parsing the CSV file,\n Error: %s", err)
			return csvRows
		}
		csvRows = append(csvRows, row)
	}
	return csvRows
}

// Creates a map of CSV headers and values
func CSVRowsToMap(header []string, values [][]string) []map[string]interface{} {
	var entities []map[string]interface{}

	for _, row := range values {
		entry := map[string]interface{}{}
		for i, value := range row {
			headerAttribute := header[i]
			// Split CSV header attributes for nested attributes
			attributeSlice := strings.Split(headerAttribute, ".")
			internal := entry
			for index, val := range attributeSlice {
				// split csv header attributes to keys
				key, attributeIndex := AttributesIndex(val)
				if attributeIndex != -1 {
					if internal[key] == nil {
						internal[key] = []interface{}{}
					}
					internalAttributeArray := internal[key].([]interface{})
					if index == len(attributeSlice)-1 {
						internalAttributeArray = append(internalAttributeArray, value)
						internal[key] = internalAttributeArray
						break
					}
					if attributeIndex >= len(internalAttributeArray) {
						internalAttributeArray = append(internalAttributeArray, map[string]interface{}{})
					}
					internal[key] = internalAttributeArray
					internal = internalAttributeArray[attributeIndex].(map[string]interface{})
				} else {
					if index == len(attributeSlice)-1 {
						internal[key] = value
						break
					}
					if internal[key] == nil {
						internal[key] = map[string]interface{}{}
					}
					internal = internal[key].(map[string]interface{})
				}
			}
		}
		entities = append(entities, entry)
	}
	return entities
}

// Recives map of entities and returns json formated result
func MapToString(entries []map[string]interface{}) (string, error) {
	//MarshalIndent for json formating output
	bytes, err := json.MarshalIndent(entries, "", "	")
	if err != nil {
		fmt.Printf("Something went wrong while marshalling byte slice, %s\n", err)
		return "", errors.New(err.Error())
	}
	return string(bytes), nil
}

// Finds attribute and index of nested attributes and returns them
func AttributesIndex(attribute string) (string, int) {
	i := strings.Index(attribute, "[")
	if i >= 0 {
		j := strings.Index(attribute, "]")
		if j >= 0 {
			index, _ := strconv.Atoi(attribute[i+1 : j])
			return attribute[0:i], index
		}
	}
	return attribute, -1
}

// Creates json file
func WriteToJSON(json string) {
	// 0777 for read, write, & execute for owner, group and others
	err := os.WriteFile("./result.json", []byte(json), 0777)
	CheckError(err)
}

// Checking errors
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
