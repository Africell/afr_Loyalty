package Lendme

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func (Uc *UserControl) Import_Subscribers_ARPU_LineByLine(FileName string) (err error) {
	log.Println("***Subscribers ARPU importing process started ***")
	//Section 1: open file for reading
	if FileName == "" {
		FileName = Configuration.ARPU_File_Path
	}
	file, err := os.Open(FileName)
	if err != nil {
		log.Println("Error opening subscribers ARPU file " + FileName + ": " + err.Error())
		return err
	}
	defer file.Close()

	// Section 2: read lines
	r := bufio.NewReader(file)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Println("***Subscribers ARPU importing process finished ***")
				return nil
			} else {
				log.Println("Error reading line from subscribers ARPU file: " + err.Error())
				return err
			}
		}
		if len(line) > 0 {
			log.Println(line)
		}
	}
}

func (Uc *UserControl) Import_Subscribers_ARPU_bychunks() {
	log.Println("***Subscribers ARPU importing processstarted ***")
	var chunkSize = 10
	f, err := os.Open(Configuration.ARPU_File_Path)
	if err != nil {
		log.Println("Failed to open subscribers ARPU file")
		return
	}
	// remember to close the file at the end of the program
	defer f.Close()

	buf := make([]byte, chunkSize)

	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if err == io.EOF {
			break
		}

		fmt.Println(string(buf[:n]))
	}

	log.Println("***Subscribers ARPU importing process finished ***")
}
