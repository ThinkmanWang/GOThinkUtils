package thinkutils

import (
	"bufio"
	"os"
)

type fileutils struct {
}

type OnReadLineCallback func(nLine uint32, szLine string)

func (this fileutils) ReadLine(szPath string, callback OnReadLineCallback) {
	inFile, err := os.Open(szPath)
	if err != nil {
		return
	}

	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	var nLine uint32 = 0
	for scanner.Scan() {
		if callback != nil {
			callback(nLine, scanner.Text())
		}
		nLine++
		//fmt.Println(scanner.Text()) // the line
	}
}
