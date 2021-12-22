package thinkutils

import (
	"bufio"
	"errors"
	"io"
	"os"
)

type fileutils struct {
}

type OnReadLineCallback func(nLine uint32, szLine string)

func (this fileutils) Exists(szPath string) bool {
	_, err := os.Stat(szPath)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func (this fileutils) MkDir(szPath string) {
	if _, serr := os.Stat(szPath); serr != nil {
		merr := os.MkdirAll(szPath, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}

func (this fileutils) ReadLine(szPath string, callback OnReadLineCallback) {
	inFile, err := os.Open(szPath)
	if err != nil {
		return
	}

	defer inFile.Close()

	var nLine uint32 = 0
	reader := bufio.NewReader(inFile)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		szLine := string(lineBytes)
		if callback != nil {
			callback(nLine, szLine)
		}
		nLine++
		//utfStr := ConvertEncoding(gbkStr, "GBK")
		//fmt.Println(utfStr)
	}

	//scanner := bufio.NewScanner(inFile)
	//var nLine uint32 = 0
	//for scanner.Scan() {
	//	if callback != nil {
	//		callback(nLine, scanner.Text())
	//	}
	//	nLine++
	//	//fmt.Println(scanner.Text()) // the line
	//}
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
