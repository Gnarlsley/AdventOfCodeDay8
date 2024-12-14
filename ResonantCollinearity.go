package main

import (
	"bufio"
	"fmt"
	"os"
)

const puzzleText = "test.txt"

// storing coordinates in a matrix for matrix multiplication purposes
type Matrix struct {
	Rows  int
	Cols  int
	data  []int
	value rune
}

// Radio towers have a symbol and a list of coordinates
type RadioTowers struct {
	Char      rune
	Locations []*Matrix
}

// default matrix constructor
func newMatrix(rows, cols int) *Matrix {
	return &Matrix{
		Rows: rows,
		Cols: cols,
		data: make([]int, rows*cols),
	}
}

// matrix constructor with a symbol attached
func newMatrixValuePair(rows, cols int, value rune) *Matrix {
	return &Matrix{
		Rows:  rows,
		Cols:  cols,
		data:  make([]int, rows*cols),
		value: value,
	}
}

// getter and setter for matrix data
func (m *Matrix) At(row, col int) int {
	return m.data[row*m.Cols+col]
}

func (m *Matrix) Set(row, col, value int) {
	m.data[row*m.Cols+col] = value
}

// multiply a 2x1 matrix by a 180 rotation matrix use like rotatedMatrix := matrix.ApplyRotation()
func (n *Matrix) ApplyRotation() *Matrix {
	c_11 := -1 * n.At(0, 0)
	c_21 := -1 * n.At(1, 0)

	matrix := newMatrix(2, 1)
	matrix.Set(0, 0, c_11)
	matrix.Set(1, 0, c_21)
	return matrix
}

// constructor for radio towers
func newRadioTowerEntry(character rune, location []*Matrix) *RadioTowers {
	return &RadioTowers{
		Char:      character,
		Locations: location,
	}
}

// utility function to determine if a rune is present in a rune slice
func NotInSlice(element rune, slice []rune) bool {
	for _, item := range slice {
		if element == item {
			return false
		}
	}
	return true
}

// read data from file
func GetData() []string {

	fileName := puzzleText
	var data []string
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return data
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}

	return data
}

// parse data from string to runes
func ParseData(data []string) []rune {
	var runes []rune

	for _, element := range data {
		runes = append(runes, []rune(element)...)
	}
	return runes
}

// get file dimensions
func GetFileInfo() ([]rune, int, int) {
	rawData := GetData()
	parsedData := ParseData(rawData)
	rows := len(rawData)
	cols := len(parsedData) / rows
	return parsedData, rows, cols
}

// convert rune slice to matrix slice storing coordinate points and tower symbol
func DataMatrix() []*Matrix {
	var coords []*Matrix

	data, rows, cols := GetFileInfo()

	for j := 0; j < rows; j++ {
		for i := 0; i < cols; i++ {
			element := data[j*cols+i]
			matrix := newMatrixValuePair(2, 1, element)
			matrix.Set(0, 0, j)
			matrix.Set(1, 0, i)
			coords = append(coords, matrix)

		}
	}

	// for idx, matrix := range coords {
	// 	fmt.Printf("Matrix %d:\n", idx+1)
	// 	for i := 0; i < matrix.Rows; i++ {
	// 		for j := 0; j < matrix.Cols; j++ {
	// 			fmt.Printf("%d ", matrix.data[i*matrix.Cols+j])
	// 		}
	// 		fmt.Println() // Newline after each row
	// 	}
	// 	fmt.Println()
	// }

	return coords

}

// sort matrix slice by tower symbol
func LocateRadioTowers() []*RadioTowers {
	coords := DataMatrix()
	var TowerTypes []rune
	var RadioTowersLocations []*RadioTowers
	for _, matrix := range coords {
		if NotInSlice(matrix.value, TowerTypes) {
			TowerTypes = append(TowerTypes, matrix.value)
		}
	}

	for _, Tower := range TowerTypes {
		var tempLocations []*Matrix
		for _, matrix := range coords {
			if matrix.value == Tower {
				tempLocations = append(tempLocations, matrix)
			}
		}
		RadioTowerLocation := newRadioTowerEntry(Tower, tempLocations)
		RadioTowersLocations = append(RadioTowersLocations, RadioTowerLocation)
	}

	return RadioTowersLocations
}

func main() {
	data := LocateRadioTowers()
	for _, entry := range data {
		println(string(entry.Char))
	}
}
