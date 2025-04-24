// --------------------------------------------------
//
// HOMEWORK 3
//
// Due: Sunday, Mar 16, 2025 (23h59)
//
// Name: Mateo Brancoli, Ertug Umsur
//
// Email:mbrancoli@olin.edu, eumsur@olin.edu
//
// Remarks, if any:
//
// some additional help provided by Nividh
//
// had an issue with my joins where we could not account
// for the nil values in some rows, making some of our
// columns shift left. This issue is seen on our last unit
// test where we got an 'invalid groupBy field' error due
// to the shifting of columns caused by the inner join...
//
// --------------------------------------------------
//
// Please fill in this file with your solutions and submit it
//
// The functions below are stubs that you should replace with your
// own implementation.
//
// PLEASE DO NOT CHANGE THE SIGNATURE IN THE STUBS BELOW.
// Doing so makes it impossible for me to test your code.
//
// --------------------------------------------------

package main

import (
	"fmt"
	"strings"
)

type Record []string

type Table struct {
	fields     []string
	rows       []Record
	primaryKey string
}

func newTable(fields []string) Table {
	rows := make([]Record, 0, 10)
	return Table{fields, rows, ""}
}

func CreateTable(fields []string, primaryKey string) Table {
	t := newTable(fields)
	t.primaryKey = primaryKey
	return t
}

func getWidths(t Table) ([]int, int) {
	// Returns the width of the largest value in each field, in order.
	// Also return the total width.
	result := make([]int, len(t.fields))
	total := 0
	for i, f := range t.fields {
		currMax := len(f)
		for _, r := range t.rows {
			l := len(r[i])
			if l > currMax {
				currMax = l
			}
		}
		result[i] = currMax
		total += currMax
	}
	return result, total
}

func PrintTable(t Table) {
	widths, total := getWidths(t)
	line := strings.Repeat("-", total+len(t.fields)*3-1)
	fmt.Printf("+%s+\n", line)
	for i, f := range t.fields {
		fmt.Printf("| %*s ", -widths[i], f)
	}
	fmt.Printf("|\n")
	fmt.Printf("+%s+\n", line)
	for _, r := range t.rows {
		for i, _ := range t.fields {
			fmt.Printf("| %*s ", -widths[i], r[i])
		}
		fmt.Printf("|\n")
	}
	fmt.Printf("+%s+\n", line)
}

//Q1

func InsertRow(t *Table, rec Record) {

	// insert the new record
	t.rows = append(t.rows, rec)

}

//Q2

//additional functions

//checks if slice contains value

func contains(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1

}

//finds index of a field in the table

func indexField(t *Table, f string) int {

	// returns the index of a field in a table
	return contains(t.fields, f)

}

func Project(t Table, fields []string) Table {

	//new table
	new_t := CreateTable(fields, t.primaryKey)

	//go through each row and add just the fields we want
	for _, row := range t.rows {
		var new_row Record
		for _, field := range fields {
			new_row = append(new_row, row[indexField(&t, field)])
		}
		InsertRow(&new_t, new_row)
	}

	return new_t

}

//Q3

func Filter(t Table, pred func(Record) bool) Table {

	//filtering rows, make a slice
	filteredRows := make([]Record, 0)
	for _, r := range t.rows {

		//this function magically returns bools, idk what it does
		if pred(r) {

			//apend to slice filtered results
			filteredRows = append(filteredRows, r)

		}

	}

	//update rows with filtered results
	t.rows = filteredRows
	return t

}

//Q4

func Join(t1 Table, t2 Table) Table {

	//generalized initiation portion of code
	//create new table with fields of both tables
	newFields := append(t1.fields, t2.fields...)
	newTable := CreateTable(newFields, "")

	//concatenate every record from t1 with every record from t2
	for _, row1 := range t1.rows {
		for _, row2 := range t2.rows {
			newRow := append(row1, row2...)
			newTable.rows = append(newTable.rows, newRow)
		}
	}

	return newTable
}

//Q5

func InnerJoin(t1 Table, t2 Table, f1 string, f2 string) Table {

	//generalized initiation portion of code
	//create new table with fields of both tables
	newFields := append(t1.fields, t2.fields...) //the ... is necessary to expand elements of variables into individual elements, which are then appended to the first variable.
	newTable := CreateTable(newFields, "")

	//define to iterate over indecies
	index1, index2 := 0, 0
	for i, field := range t1.fields {
		if field == f1 {
			index1 = i
			break
		}
	}
	for i, field := range t2.fields {
		if field == f2 {
			index2 = i
			break
		}
	}

	// inner join
	for _, row1 := range t1.rows {
		for _, row2 := range t2.rows {
			if row1[index1] == row2[index2] {
				newRow := append(row1, row2...)
				newTable.rows = append(newTable.rows, newRow)
			}
		}
	}
	return newTable

}

func LeftOuterJoin(t1 Table, t2 Table, f1 string, f2 string) Table {

	//generalized initiation portion of code

	//define to iterate over indecies
	index1, index2 := 0, 0
	for i, field := range t1.fields {
		if field == f1 {
			index1 = i
			break
		}
	}
	for i, field := range t2.fields {
		if field == f2 {
			index2 = i
			break
		}
	}

	//create new table with fields of both tables
	newFields := append(t1.fields, t2.fields...)
	newTable := CreateTable(newFields, "")

	// left outer join (same block as before and...)
	for _, row1 := range t1.rows {
		matched := false
		for _, row2 := range t2.rows {
			if row1[index1] == row2[index2] {
				newRow := append(row1, row2...)
				newTable.rows = append(newTable.rows, newRow)
				matched = true
			}
		}

		//(add new block)new row with nil or empty values for the columns from t2
		if !matched {
			newRow := append(row1, make([]string, len(t2.fields))...)
			newTable.rows = append(newTable.rows, newRow)
		}
	}

	return newTable
}

func RightOuterJoin(t1 Table, t2 Table, f1 string, f2 string) Table {

	//generalized initiation portion of code

	//define to iterate over indecies
	index1, index2 := 0, 0
	for i, field := range t1.fields {
		if field == f1 {
			index1 = i
			break
		}
	}
	for i, field := range t2.fields {
		if field == f2 {
			index2 = i
			break
		}
	}

	//create new table with fields of both tables
	newFields := append(t1.fields, t2.fields...)
	newTable := CreateTable(newFields, "")

	// right outer join (switch 1s for 2s on both blocks)
	for _, row2 := range t2.rows {
		matched := false
		for _, row1 := range t1.rows {
			if row1[index1] == row2[index2] {
				newRow := append(row1, row2...)
				newTable.rows = append(newTable.rows, newRow)
				matched = true
			}
		}

		//new row with nil or empty values for the columns from t1
		if !matched {
			newRow := append(row2, make([]string, len(t2.fields))...)
			newTable.rows = append(newTable.rows, newRow)
		}
	}

	return newTable
}

func FullOuterJoin(t1 Table, t2 Table, f1 string, f2 string) Table {

	//generalized initiation portion of code

	//define to iterate over indecies
	index1, index2 := 0, 0
	for i, field := range t1.fields {
		if field == f1 {
			index1 = i
			break
		}
	}
	for i, field := range t2.fields {
		if field == f2 {
			index2 = i
			break
		}
	}

	//create new table with fields of both tables
	newFields := append(t1.fields, t2.fields...)
	newTable := CreateTable(newFields, "")

	// full outer join (add left and right outer join blocks)

	//left outer
	for _, row1 := range t1.rows {
		matched := false
		for _, row2 := range t2.rows {
			if row1[index1] == row2[index2] {
				newRow := append(row1, row2...)
				newTable.rows = append(newTable.rows, newRow)
				matched = true
			}
		}

		//(add new block)new row with nil or empty values for the columns from t2
		if !matched {
			newRow := append(row1, make([]string, len(t2.fields))...)
			newTable.rows = append(newTable.rows, newRow)
		}
	}

	//right outer
	for _, row2 := range t2.rows {
		matched := false
		for _, row1 := range t1.rows {
			if row1[index1] == row2[index2] {
				newRow := append(row1, row2...)
				newTable.rows = append(newTable.rows, newRow)
				matched = true
			}
		}

		//new row with nil or empty values for the columns from t1
		if !matched {
			newRow := append(row2, make([]string, len(t2.fields))...)
			newTable.rows = append(newTable.rows, newRow)
		}
	}

	return newTable
}

//Q6

func Aggregate(t Table, groupBy string, concat []string) Table {

	//find indecies of the groupBy and concat fields
	groupByIndex := indexField(&t, groupBy)
	if groupByIndex == -1 {
		fmt.Println("Invalid groupBy field")
		return CreateTable([]string{"dummy"}, "")
	}
	concatIndices := make([]int, len(concat))
	for i, field := range concat {
		concatIndices[i] = indexField(&t, field)
		if concatIndices[i] == -1 {
			fmt.Println("Invalid concat field:", field)
			return CreateTable([]string{"dummy"}, "")
		}
	}

	//map to store the aggregated results
	groupedRecords := make(map[string]Record)
	for _, row := range t.rows {
		groupKey := row[groupByIndex]
		if _, exists := groupedRecords[groupKey]; !exists {
			groupedRecords[groupKey] = make(Record, len(concat))
		}
		for i, index := range concatIndices {
			groupedRecords[groupKey][i] += row[index] + " "
		}
	}

	//new table with the groupBy field and the concat fields
	newFields := append([]string{groupBy}, concat...)
	newTable := CreateTable(newFields, "")

	// insert aggregated records into new table
	for groupKey, aggregatedRecord := range groupedRecords {
		newRow := append([]string{groupKey}, aggregatedRecord...)
		InsertRow(&newTable, newRow)
	}

	return newTable
}

func main() {
	// Sample outputs.

	//Q1
	t := CreateTable([]string{"n", "square"}, "n")
	for i := 0; i < 20; i++ {
		InsertRow(&t, Record{fmt.Sprintf("%d", i), fmt.Sprintf("%d", i*i)})
	}
	PrintTable(t)

	//Q2

	PrintTable(Project(t, []string{"square"}))

	//Q3

	PrintTable(Filter(t, func(r Record) bool { return strings.HasPrefix(r[0], "1") }))
	t2 := CreateTable([]string{"number", "text"}, "number")
	for i := 0; i < 15; i++ {
		InsertRow(&t2, Record{fmt.Sprintf("%d", i*2), fmt.Sprintf("value_%d", i*2)})
	}
	PrintTable(t2)

	//Q4

	PrintTable(Join(t, t2))
	PrintTable(Filter(Join(t, t2), func(r Record) bool { return r[0] == r[2] }))

	//Q5

	PrintTable(InnerJoin(t, t2, "n", "number"))
	PrintTable(LeftOuterJoin(t, t2, "n", "number"))
	PrintTable(RightOuterJoin(t, t2, "n", "number"))
	PrintTable(FullOuterJoin(t, t2, "n", "number"))
	PrintTable(LeftOuterJoin(FullOuterJoin(t, t2, "n", "number"), t, "2.number", "square"))

	//Q6
	t3 := CreateTable([]string{"n", "first", "last"}, "n")
	for i := 0; i < 20; i++ {
		v := fmt.Sprintf("%d", i)
		InsertRow(&t3, Record{v, v[:1], v[len(v)-1:]})
	}
	PrintTable(t3)
	PrintTable(Aggregate(t3, "first", []string{"n", "last"}))
	PrintTable(Aggregate(t3, "last", []string{"n", "first"}))
	PrintTable(Aggregate(InnerJoin(t, t3, "n", "n"), "2.first", []string{"1.square", "1.n"}))

}
