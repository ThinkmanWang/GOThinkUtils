package thinkutils

import (
	"bufio"
	"os"
)

type fileutils struct {
}

type OnReadLineCallback func(szLine string)

func (this fileutils) ReadLine(szPath string, callback OnReadLineCallback) {
	inFile, err := os.Open(szPath)
	if err != nil {
		return
	}

	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		if callback != nil {
			callback(scanner.Text())
		}
		//fmt.Println(scanner.Text()) // the line
	}
}
