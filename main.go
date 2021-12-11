package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var log *logger.LocalLogger = logger.DefaultLogger()

func allMainFile() []string {
	szCmd := "grep \"func main()\"" + " `" + "grep \"package main\" . -rl" + "` " + "-rl"
	out, err := exec.Command("bash", "-c", szCmd).Output()
	if err != nil {
		log.Error("%s", err.Error())
	}

	szOutput := string(out)
	lstFiles := strings.Split(szOutput, "\n")

	return lstFiles
}

func buildFile(szDir string, szFile string) error {
	thinkutils.FileUtils.MkDir(szDir)

	lstItem := strings.Split(szFile, "/")
	if nil == lstItem || len(lstItem) <= 0 {
		return nil
	}

	szFileName := strings.Split(lstItem[len(lstItem)-1], ".")[0]
	//log.Info("%s => [%s]", szFile, szFileName)

	szOutput := fmt.Sprintf("%s/%s", szDir, szFileName)
	if thinkutils.FileUtils.Exists(szOutput) {
		nNum := 1
		for {
			szTemp := fmt.Sprintf("%s%d", szOutput, nNum)

			if false == thinkutils.FileUtils.Exists(szTemp) {
				szOutput = szTemp
				break
			}

			nNum++
		}
	}

	szCmd := fmt.Sprintf("go build -o %s %s", szOutput, szFile)
	log.Info(szCmd)

	_, err := exec.Command("bash", "-c", szCmd).Output()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	szPath, err := os.Getwd()
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info(szPath)

	lstFiles := allMainFile()
	if nil == lstFiles || len(lstFiles) <= 0 {
		log.Info("No file found return!")
		return
	}

	for i := 0; i < len(lstFiles); i++ {
		szFile := lstFiles[i]
		if thinkutils.StringUtils.IsEmpty(szFile) {
			continue
		}

		if false == strings.HasSuffix(szFile, ".go") {
			continue
		}

		buildFile("bin", szFile)
	}
}
