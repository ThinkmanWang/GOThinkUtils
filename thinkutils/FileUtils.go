package thinkutils

import (
	"bufio"
	"io"
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

func (this fileutils) Copy(szSrc string, szDst string) error {
	srcInfo, err := os.Stat(szSrc)
	if err != nil {
		return err
	}

	srcFile, err := os.Open(szSrc)
	if err != nil {
		return err
	}
	defer func() {
		err := srcFile.Close()
		if err != nil {
			return
		}
	}()

	destFile, err := os.Create(szDst)
	if err != nil {
		return err
	}
	defer func() {
		err := destFile.Close()
		if err != nil {
			return
		}
	}()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	if err := os.Chmod(szDst, srcInfo.Mode()); err != nil {
		return err
	}

	return destFile.Sync()
}
