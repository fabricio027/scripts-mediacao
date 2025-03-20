package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const DirCollFTTX string = "/mediacao/oi/infraco/move/fttx/bruto"

// const DirCollFTTX string = "C:\\Users\\oi414774\\Desktop\\Fabricio\\scripts\\go\\banco"
const fileListFTTX string = "/tmp/filesFTTX.tmp"

// const fileListFTTX string = "scripts\\filesFTTX.tmp"

var ListFilesFTTX []string = generateList()
var FilesWithoutColl []string

func main() {
	ListFilesFTTX = removeDuplicate()
	FilesWithoutColl = verifyFilesWithoutColl()

	fmt.Println("== VERIFICANDO STATUS COLETA FTTX ==")
	fmt.Println()
	if len(FilesWithoutColl) == 0 {
		fmt.Println("Sem arquivos FTTX em atraso.")
	} else {
		fmt.Println("Segue(m) arquivo(s) pendente(s) de coleta:")
		print()
	}

	writeFile()
	fmt.Println()
	fmt.Println("== Fim do programa ==")

}

func print() {
	for _, line := range FilesWithoutColl {
		fmt.Println(line)
	}
}

func verifyFilesWithoutColl() []string {
	result := make([]string, 0)

	for _, file := range ListFilesFTTX {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", DirCollFTTX, file)); err != nil {
			result = append(result, file)
		}
	}

	return result
}

func removeDuplicate() []string {
	result := make([]string, 0)

	for _, line := range ListFilesFTTX {
		duplicate := false
		for _, unique := range result {
			if line == unique {
				duplicate = true
				break
			}
		}

		if !duplicate {
			result = append(result, line)
		}
	}

	return result

}

func writeFile() {
	file, err := os.Create(fileListFTTX)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, line := range FilesWithoutColl {
		_, err := file.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil {
			panic(err)
		}
	}
}

func generateList() []string {
	list := make([]string, 0)
	list = append(list, generateListOld(fileListFTTX)...)
	list = append(list, generateListRecent()...)

	return list
}

func generateListOld(filesFTTX string) []string {
	if _, err := os.Stat(filesFTTX); err != nil {
		file, err := os.Create(fileListFTTX)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	file, err := os.Open(filesFTTX)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lines := make([]string, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return lines
}

func generateListRecent() []string {
	date := time.Now().Add(-24 * time.Hour)
	list := []string{
		fmt.Sprintf("ALV-RJ-NACF-01-%s.gz@INFRACO#ALV-RJ-NACF-01-%s.gz", date.Format("20060102"), date.Format("20060102")),
		fmt.Sprintf("ALV-RJ-NACF-02-%s.gz@INFRACO#ALV-RJ-NACF-02-%s.gz", date.Format("20060102"), date.Format("20060102")),
		fmt.Sprintf("BDEA-BA-NACF-02-%s.gz@INFRACO#BDEA-BA-NACF-02-%s.gz", date.Format("20060102"), date.Format("20060102")),
		fmt.Sprintf("CTBH-PR-NACF-02-%s.gz@INFRACO#CTBH-PR-NACF-02-%s.gz", date.Format("20060102"), date.Format("20060102")),
		fmt.Sprintf("CTME-PR-NACF-01-%s.gz@INFRACO#CTME-PR-NACF-01-%s.gz", date.Format("20060102"), date.Format("20060102")),
		fmt.Sprintf("CTME-PR-NACF-02-%s.gz@INFRACO#CTME-PR-NACF-02-%s.gz", date.Format("20060102"), date.Format("20060102")),
		fmt.Sprintf("PRA-RJ-NACF-01-%s.gz@INFRACO#PRA-RJ-NACF-01-%s.gz", date.Format("20060102"), date.Format("20060102")),
		fmt.Sprintf("BDEA-BA-NACF-01-%s.gz@INFRACO#BDEA-BA-NACF-01-%s.gz", date.Format("20060102"), date.Format("20060102")),
		fmt.Sprintf("CTBH-PR-NACF-01-%s.gz@INFRACO#CTBH-PR-NACF-01-%s.gz", date.Format("20060102"), date.Format("20060102")),
		fmt.Sprintf("PRA-RJ-NACF-02-%s.gz@INFRACO#PRA-RJ-NACF-02-%s.gz", date.Format("20060102"), date.Format("20060102")),
		fmt.Sprintf("BRRR-BA-NACF-01-%s.gz@INFRACO#BRRR-BA-NACF-01-%s.gz", date.Format("20060102"), date.Format("20060102")),
		fmt.Sprintf("BRRR-BA-NACF-02-%s.gz@INFRACO#BRRR-BA-NACF-02-%s.gz", date.Format("20060102"), date.Format("20060102")),
	}

	return list
}
