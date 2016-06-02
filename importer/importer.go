package importer

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Just to initialize MySQL Driver
)

// Import IPTU data
func Import(filename string) {

	db := connectDb()
	defer db.Close()

	file := openFile(filename)
	defer file.Close()

	print("Truncating table...")
	_, err := db.Exec("TRUNCATE TABLE `iptu`")
	if err != nil {
		panic(err.Error())
	}
	println(" Table truncated!")

	stmtIns, err := db.Prepare("INSERT INTO `iptu` VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if i > 1000000 {
			if i%1000 == 0 {
				fmt.Printf("Importing %s\n", RenderInteger("#,###.", i))
			}
			slices := strings.Split(scanner.Text(), "|")

			// Data Cadastramento
			slices[3] = fmt.Sprintf("20%s-%s-%s", slices[3][6:8], slices[3][3:5], slices[3][0:2])

			// Tipo 1
			if strings.Contains(slices[4], "CPF") {
				slices[4] = "CPF"
				slices[5] = slices[5][2 : len(slices[5])-1]
			} else if strings.Contains(slices[4], "CNPJ") {
				slices[4] = "CNPJ"
			} else {
				slices[4] = ""
			}
			// Tipo 2
			if strings.Contains(slices[7], "CPF") {
				slices[7] = "CPF"
				slices[8] = slices[8][2 : len(slices[8])-1]
			} else if strings.Contains(slices[7], "CNPJ") {
				slices[7] = "CNPJ"
			} else {
				slices[7] = ""
			}

			// Nome 1
			if len(slices[6]) > 1 && (slices[6][0] == '\'' || slices[6][0] == ',' || slices[6][0] == '.') {
				slices[6] = slices[6][1 : len(slices[6])-1]
			}
			// Nome 2
			if len(slices[9]) > 1 && (slices[9][0] == '\'' || slices[9][0] == ',' || slices[9][0] == '.') {
				slices[9] = slices[9][1 : len(slices[9])-1]
			}

			// Numericos
			slices[19] = convertNum(slices[19])
			slices[23] = convertNum(slices[23])
			slices[24] = convertNum(slices[24])
			slices[27] = convertNum(slices[27])
			slices[31] = convertNum(slices[31])

			// Inteiros
			slices[18] = convertInt(slices[18])
			slices[20] = convertInt(slices[20])
			slices[21] = convertInt(slices[21])
			slices[22] = convertInt(slices[22])
			slices[26] = convertInt(slices[26])

			//Ano Construção
			if slices[25] == "0" {
				slices[25] = ""
			}

			_, err = stmtIns.Exec(*(convertSlice(slices))...)
			if err != nil {
				panic(err.Error())
			}
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func openFile(filename string) *os.File {
	print("Opening file...")
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	println(" File opened!")
	return file
}

func connectDb() *sql.DB {
	print("Connecting to DB... ")
	db, err := sql.Open("mysql", "iptu:iptu@/iptu")
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	println(" Connected to DB!")
	return db
}

func convertSlice(slice []string) *[]interface{} {
	result := make([]interface{}, len(slice))
	for i := range slice {
		result[i] = slice[i]
	}
	return &result
}

func convertNum(input string) string {
	return strings.Replace(input, ",", ".", 1)
}

func convertInt(input string) string {
	return strings.Replace(input, ",", ".", 1)
}
