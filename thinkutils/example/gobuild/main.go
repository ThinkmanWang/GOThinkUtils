package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	log         *logger.LocalLogger = logger.DefaultLogger()
	g_nTotal    int                 = 0
	g_nPos      int                 = 0
	g_nPlatform int                 = 3
)

func allMainFile() []string {
	szCmd := "grep \"func main()\"" + " `" + "grep \"package main\" . -rl" + "` " + "-rl"
	out, err := exec.Command("bash", "-c", szCmd).Output()
	if err != nil {
		log.Error("%s", err.Error())
	}

	szOutput := string(out)
	lstFiles := strings.Split(szOutput, "\n")

	lstRet := make([]string, 0)
	for i := 0; i < len(lstFiles); i++ {
		if thinkutils.StringUtils.IsEmpty(lstFiles[i]) {
			continue
		}

		if false == strings.HasSuffix(lstFiles[i], ".go") {
			continue
		}

		lstRet = append(lstRet, lstFiles[i])
	}

	return lstRet
}

func buildFile(szEnv string, szDir string, szFile string, szExt string) error {
	g_nPos++
	thinkutils.FileUtils.MkDir(szDir)

	lstItem := strings.Split(szFile, "/")
	if nil == lstItem || len(lstItem) <= 0 {
		return nil
	}

	szFileName := strings.Split(lstItem[len(lstItem)-1], ".")[0]
	//log.Info("%s => [%s]", szFile, szFileName)

	szOutput := fmt.Sprintf("%s/%s", szDir, szFileName)
	if thinkutils.FileUtils.Exists(fmt.Sprintf("%s%s", szOutput, szExt)) {
		nNum := 1
		for {
			szTemp := fmt.Sprintf("%s%d", szOutput, nNum)

			if false == thinkutils.FileUtils.Exists(fmt.Sprintf("%s%s", szTemp, szExt)) {
				szOutput = szTemp
				break
			}

			nNum++
		}
	}

	szCmd := fmt.Sprintf("%s go build -o %s%s %s", szEnv, szOutput, szExt, szFile)
	log.Info("[%d/%d] %s", g_nPos, g_nTotal, szCmd)

	out, err := exec.Command("bash", "-c", szCmd).Output()
	if err != nil {
		fmt.Println(strings.TrimSpace(string(out)))
		return err
	}

	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

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

	g_nTotal = len(lstFiles) * g_nPlatform
	for i := 0; i < len(lstFiles); i++ {
		buildFile("env GOOS=linux GOARCH=amd64", "bin/linux", lstFiles[i], "")
		buildFile("env GOOS=darwin GOARCH=amd64", "bin/mac", lstFiles[i], "")
		buildFile("env GOOS=windows GOARCH=amd64", "bin/windows", lstFiles[i], ".exe")
	}
}
