package main

import (
	"bufio"
	"fmt"
	"os"
)

const puzzleText = "puzzle.txt"

//417 too high
//377 too low

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

func (m *Matrix) Equals(other *Matrix) bool {
	if m == nil || other == nil {
		return false
	}

	for i := range m.data {
		if m.data[i] != other.data[i] {
			return false
		}
	}

	return true
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

func MatrixNotInSlice(element *Matrix, slice []*Matrix) bool {
	for _, item := range slice {
		if element.data[0] == item.data[0] && element.data[1] == item.data[1] {
			return false
		}
	}
	return true
}

func computeRelativePosition(origin, other *Matrix) *Matrix {
	tempMatrix := newMatrix(2, 1)
	fromOriginX := other.At(0, 0) - origin.At(0, 0)
	fromOriginY := other.At(1, 0) - origin.At(1, 0)
	tempMatrix.Set(0, 0, fromOriginX)
	tempMatrix.Set(1, 0, fromOriginY)
	return tempMatrix.ApplyRotation()
}

func computeReflectedPosition(rotatedMatrix, origin *Matrix, rows, cols int) *Matrix {
	rotatedFromOriginX := rotatedMatrix.At(0, 0)
	rotatedFromOriginY := rotatedMatrix.At(1, 0)

	reflectedX := rotatedFromOriginX + origin.At(0, 0)
	reflectedY := rotatedFromOriginY + origin.At(1, 0)

	if reflectedX >= cols || reflectedX < 0 || reflectedY >= rows || reflectedY < 0 {
		return nil
	}

	result := newMatrix(2, 1)
	result.Set(0, 0, reflectedX)
	result.Set(1, 0, reflectedY)
	result.value = '#'
	return result
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
func GetFileInfo() (int, int) {
	rawData := GetData()
	parsedData := ParseData(rawData)
	rows := len(rawData)
	cols := len(parsedData) / rows
	return rows, cols
}

// convert rune slice to matrix slice storing coordinate points and tower symbol
func DataMatrix() []*Matrix {
	var coords []*Matrix

	rawData := GetData()
	data := ParseData(rawData)
	rows, cols := GetFileInfo()

	for j := 0; j < rows; j++ {
		for i := 0; i < cols; i++ {
			element := data[j*cols+i]
			matrix := newMatrixValuePair(2, 1, element)
			matrix.Set(0, 0, i)
			matrix.Set(1, 0, j)
			coords = append(coords, matrix)

		}
	}

	// for idx, matrix := range coords {
	// 	fmt.Printf("Matrix %d:\n", idx+1)
	// 	for i := 0; i < matrix.Rows; i++ {
	// 		for j := 0; j < matrix.Cols; j++ {
	// 			fmt.Printf("%d ", matrix.data[i*matrix.Cols+j])
	// 		}
	// 		fmt.Println()
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
	rows, cols := GetFileInfo()
	var TowerRangeLocations []*Matrix
	for _, entry := range data {
		symbol := entry.Char
		locations := entry.Locations
		if symbol == '.' {
			continue
		}
		for i, origin := range locations {
			for j, other := range locations {
				if i == j {
					continue
				}
				rotatedMatrix := computeRelativePosition(origin, other)
				rotatedMatrix = computeReflectedPosition(rotatedMatrix, origin, rows, cols)
				if rotatedMatrix != nil && MatrixNotInSlice(rotatedMatrix, TowerRangeLocations) {
					TowerRangeLocations = append(TowerRangeLocations, rotatedMatrix)
				}
			}
		}
	}

	// for _, entry := range data {
	// 	for _, RadioLocation := range entry.Locations {
	// 		for index, SignalLocation := range TowerRangeLocations {
	// 			if RadioLocation.Equals(SignalLocation) {
	// 				//remove the duplicate
	// 				if RadioLocation.value != '.' {
	// 					TowerRangeLocations = append(TowerRangeLocations[:index], TowerRangeLocations[index+1:]...)
	// 					break
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// for _, entry := range TowerRangeLocations {
	// 	println(string(entry.value))
	// 	println("X:", entry.data[0])
	// 	println("Y:", entry.data[1])
	// }
	// println(len(TowerRangeLocations))

	//TowerRangeLocations includes all the '#'
	//data includes all the towers symbols and their locations

	matrixData := DataMatrix()
	for index, location := range matrixData {
		found := false
		for _, towerLocation := range TowerRangeLocations {
			if location.Equals(towerLocation) {
				print(string(towerLocation.value))
				found = true
				break
			}
		}

		if !found {
			print(string(location.value))
		}
		if (index+1)%cols == 0 {
			print("\n")
		}
	}

	println(len(TowerRangeLocations))
}
