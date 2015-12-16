package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	baseEnv       = "ZOOKEEPER_"
	baseEnvSuffix = "_SERVICE_HOST"
)

func main() {
	var (
		envs        []string
		filePath    string
		dataToWrite string
		fileName    string
		dataDir     string
	)

	flag.StringVar(&filePath, "filepath", "", "File path for zoo.cfg")
	flag.StringVar(&fileName, "filename", "", "Filename for server configs.")
	flag.StringVar(&dataDir, "dataDir", "", "dataDir location for Zookeeper.")
	flag.Parse()

	envVars := genEnvVars()
	for idx, envVar := range envVars {
		if env, found := os.LookupEnv(envVar); found {
			envs = append(envs, fmt.Sprintf("server.%d=%s:2888:3888", idx+1, env))
		}
	}
	for idx, env := range envs {
		dataToWrite = dataToWrite + fmt.Sprintf("%s", env)
		if idx != len(envs)-1 {
			dataToWrite = dataToWrite + fmt.Sprintf("\n")
		}
	}

	file, err := os.OpenFile(filePath+fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		fmt.Printf("no file written.\n")
		return
	}
	defer file.Close()
	_, err = file.Write([]byte(dataToWrite))
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	genMyID(dataDir)

}

func genEnvVars() []string {
	limit := 99
	var envVars []string
	for ii := 1; ii <= limit; ii++ {
		envVar := baseEnv + fmt.Sprintf("%02d", ii) + baseEnvSuffix
		envVars = append(envVars, envVar)
	}
	return envVars
}

func genMyID(dataDir string) {
	var myID string
	if len(dataDir) > 1 {
		dataDir = "/tmp/zookeeper"
	}
	myID = os.Getenv("MYID")

	os.Mkdir(fmt.Sprintf("%s/%s", dataDir, myID), 0777)
	fmt.Printf("writing myid file to %s/\n", fmt.Sprintf("%s/%s", dataDir, myID))
	idFile := fmt.Sprintf("%s/myid", dataDir)
	file, err := os.OpenFile(idFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	defer file.Close()
	if err != nil {
		fmt.Printf("error opening file: %s\n", err)
		fmt.Println("no file written.")
	}
	_, err = file.Write([]byte(myID))
	if err != nil {
		fmt.Printf("error writing myID: %s", err)
	}
}
