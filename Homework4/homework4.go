// --------------------------------------------------
//
// HOMEWORK 4
//
// Due: Sunday, Apr 20, 2025 (23h59)
//
// Name: Mateo Bracoli
//
// Email: mbrancoli@olin.edu
//
// Remarks, if any:
// Teaming with Ertug
// Forgot to submit on my side
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
	"strconv"
	"strings"
)

const recordsPerBlock = 5

type Block = [recordsPerBlock]string

type File struct {
	blocks   []Block
	overflow *Block
}

func NewBlock() Block {
	return [recordsPerBlock]string{}
}

func NewSeqFile() *File {
	return &File{make([]Block, 0), nil}
}

func NewOrdFile() *File {
	b := NewBlock()
	return &File{make([]Block, 0), &b}
}

func ReadBlock(file *File, idx int) Block {
	if idx < 0 {
		fmt.Printf("[Reading overflow block]\n")
		return *file.overflow
	}
	fmt.Printf("[Reading block %d]\n", idx)
	return file.blocks[idx]
}

func WriteBlock(file *File, idx int, block Block) {
	if idx < 0 {
		fmt.Printf("[Writing overflow block]\n")
		*file.overflow = block
		return
	}
	fmt.Printf("[Writing block %d]\n", idx)
	file.blocks[idx] = block
}

func CreateBlock(file *File) int {
	b := NewBlock()
	file.blocks = append(file.blocks, b)
	return FileSize(file) - 1
}

func FileSize(file *File) int {
	return len(file.blocks)
}

func PrintFile(file *File) {
	fmt.Println("-- BLOCKS ------------------------------")
	for j, b := range file.blocks {
		fmt.Printf("Block %d\n", j)
		for i, r := range b {
			fmt.Printf(" %2d %v\n", i, r)
		}
	}
	if file.overflow != nil {
		fmt.Println("-- OVERFLOW ----------------------------")
		for i, r := range file.overflow {
			fmt.Printf(" %2d %v\n", i, r)
		}
	}
}

// Question 1

func AppendRecord(block Block, rec string) (Block, bool) {

	if rec == "" {
		panic("Can't write empty string!")
	}

	for i, r := range block {
		if r == "" {
			block[i] = rec
			return block, true
		}
	}

	return block, false
}

func FirstRecord(block Block) string {
	return block[0]
}

func LastRecord(block Block) string {
	return block[len(block)-1]
}

func FindRecord(block Block, rec string) bool {
	for _, r := range block {
		if r == rec {
			return true
		}
	}

	return false
}

func FreeSize(block Block) int {
	empty_slots := 0
	for _, r := range block {
		if r == "" {
			empty_slots = empty_slots + 1
		}
	}

	return empty_slots
}

// Question 2: sequential files

func SF_Find(file *File, record string) bool {
	for _, r := range file.blocks {
		if FindRecord(r, record) {
			return true
		}
	}

	return false
}

func SF_Insert(file *File, record string) {
	if !SF_Find(file, record) {

		if FileSize(file) == 0 {
			CreateBlock(file)
		}

		index := FileSize(file) - 1
		block := ReadBlock(file, index)
		if FreeSize(block) > 0 {
			block, _ = AppendRecord(block, record)
			WriteBlock(file, index, block)
		} else {
			indexn := CreateBlock(file)
			block, _ = AppendRecord(ReadBlock(file, indexn), record)
			WriteBlock(file, indexn, block)
		}
	}
}

// Question 3: ordered files

func OF_Find(f *File, rec string) bool {

	low := 0
	high := len(f.blocks) - 1

	for low <= high {

		//look at the middle block
		mid := (low + high) / 2

		//use functions provided
		block := ReadBlock(f, mid)

		//find first and last records in the block
		first := FirstRecord(block)
		last := LastRecord(block)

		if rec < first {

			//sirch lower half of blocks
			high = mid - 1

		} else if rec > last {

			//search higher half of blocks
			low = mid + 1
		} else {
			return FindRecord(block, rec) || FindRecord(*f.overflow, rec)
		}
	}
	return false
}

func OF_Insert(f *File, rec string) {

	if OF_Find(f, rec) {
		return
	}

	if FreeSize(*f.overflow) == 0 {

		var records []string

		if len(f.blocks) != 0 {
			for _, b := range f.blocks {
				for _, r := range b {
					records = append(records, r)
				}
			}
		}

		for _, r := range *f.overflow {
			records = append(records, r)
		}

		records = append(records, rec)

		var filtered []string
		for _, r := range records {
			if r != "" {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		for i := 0; i < len(records); i++ {
			for j := i + 1; j < len(records); j++ {

				I_part := strings.Split(records[i], "-")
				J_part := strings.Split(records[j], "-")

				if len(I_part) < 2 || len(J_part) < 2 {
					continue
				}

				I_number, errorI := strconv.Atoi(I_part[1])
				J_number, errorJ := strconv.Atoi(J_part[1])
				if errorI != nil || errorJ != nil {
					continue
				}

				if I_number > J_number {
					records[i], records[j] = records[j], records[i]
				}
			}
		}

		block_number := (len(records) + recordsPerBlock - 1) / recordsPerBlock

		f.blocks = nil

		for i := 0; i < block_number; i++ {
			new_Block := NewBlock()
			for j := 0; j < recordsPerBlock; j++ {
				index := i*recordsPerBlock + j
				if index >= len(records) {
					break
				}
				new_Block, _ = AppendRecord(new_Block, records[index])
			}

			if len(f.blocks) < i+1 {
				CreateBlock(f)
			}

			WriteBlock(f, i, new_Block)
		}

		*f.overflow = NewBlock()

	} else {
		*f.overflow, _ = AppendRecord(*f.overflow, rec)
	}
}

func main() {
	fmt.Println("=== SEQ FILE =====================================")
	f := NewSeqFile()
	for i := 1; i < 25; i++ {
		rec := fmt.Sprintf("test-%d", i)
		SF_Insert(f, rec)
	}
	PrintFile(f)
	fmt.Printf("Looking for test-14: %v\n", SF_Find(f, "test-14"))
	fmt.Printf("looking for test-99: %v\n", SF_Find(f, "test-99"))

	fmt.Println("=== ORD FILE =====================================")
	of := NewOrdFile()
	for i := 1; i < 25; i++ {
		rec := fmt.Sprintf("test-%d", i)
		OF_Insert(of, rec)
	}
	PrintFile(of)
	fmt.Printf("Looking for test-14: %v\n", OF_Find(of, "test-14"))
	fmt.Printf("looking for test-99: %v\n", OF_Find(of, "test-99"))
}
