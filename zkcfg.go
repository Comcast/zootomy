package main

import (
	"flag"
	"fmt"
	"os"
//	"strconv"
)

const (
	baseEnv       = "ZOOKEEPER_"
	baseEnvSuffix = "_SERVICE_HOST"
)

func main() {
	var (
		envs          []string
		zooCfgPath    string
		dataToWrite   string
		zooCfgName    string
		dataDir       string
		quorum        bool
		standalone    bool
		maxConn       int
		tickTime      int
		initLimit     int
		syncLimit     int
		clientPort    int
		dynamicConfig string
	)

	flag.StringVar(&zooCfgPath, "zooCfgPath", "", "File path for zoo.cfg")
	flag.StringVar(&zooCfgName, "zooCfgName", "zoo.cfg", "Filename for server configs.")
	flag.StringVar(&dataDir, "dataDir", "", "dataDir location for Zookeeper.")
	flag.BoolVar(&quorum, "quorum", true, "enables the quorum to listen on all IPs.")
	flag.BoolVar(&standalone, "standalone", false, "enables the Zookeeper node to start in standalone mode.")
	flag.IntVar(&maxConn, "maxConn", 0, "the maximum allowable connections at once. defaults to unlimited, but this will normally require a ramdisk. see readme.")
	flag.IntVar(&tickTime, "tickTime", 2000, "the number in milliseconds of each tick.")
	flag.IntVar(&initLimit, "initLimit", 10, "the number of ticks that the initial synchronization phase can take.")
	flag.IntVar(&syncLimit, "syncLimit", 5, "the number of ticks that can pass between sending a request and getting an acknowledgement.")
	flag.IntVar(&clientPort, "clientPort", 2181, "the port to which clients will connect.")
	flag.StringVar(&dynamicConfig, "dynamicConfig", "zoo.cfg.dynamic", "the path and file name for the dynamic zookeeper.")

	flag.Parse()

	var configs = make(map[string]interface{})

	configs["dataDir"] = dataDir
	configs["quorumListenOnAllIPs"] = quorum
	configs["standaloneEnabled"] = standalone
	configs["maxConn"] = maxConn
	configs["tickTime"] = tickTime
	configs["initLimit"] = initLimit
	configs["syncLimit"] = syncLimit
	configs["clientPort"] = clientPort
	configs["dynamicConfigFile"] = dynamicConfig

	envVars := genEnvVars()

	/*
		So Zookeeper needs to bind on 0.0.0.0 for the local listener.
		Nice of the documentation to talk about that.
		Info: http://stackoverflow.com/questions/30940981/zookeeper-error-cannot-open-channel-to-x-at-election-address
	*/
//	self := "0.0.0.0"
//	me, err := strconv.Atoi(os.Getenv("MYID"))
//	if err != nil {
//		fmt.Printf("cannot read my identity: %s", err)
//	}

	for idx, envVar := range envVars {
		if env, found := os.LookupEnv(envVar); found {
//			if idx+1 == me {
//				envs = append(envs, fmt.Sprintf("server.%d=%s:2888:3888", idx+1, self))
//			} else {
//				envs = append(envs, fmt.Sprintf("server.%d=%s:2888:3888", idx+1, env))
//			}
			envs = append(envs, fmt.Sprintf("server.%d=%s:2888:3888", idx+1, env))
		}
	}
	for idx, env := range envs {
		dataToWrite = dataToWrite + fmt.Sprintf("%s", env)
		if idx != len(envs)-1 {
			dataToWrite = dataToWrite + fmt.Sprintf("\n")
		}
	}

	// writes the dynamic file locations.
	file, err := os.OpenFile(dynamicConfig, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		fmt.Printf("no file written.\n")
		return
	}
	defer file.Close()
	fmt.Printf("wrote the dynamic config to: %s\n", dynamicConfig)
	_, err = file.Write([]byte(dataToWrite))
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	err = genConfig(configs, zooCfgPath+zooCfgName)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	genMyID(dataDir)

}

// genConfig writes the Zookeeper configuration to disk.
func genConfig(c map[string]interface{}, f string) error {
	file, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	defer file.Close()
	if err != nil {
		return err
	}
	for a, b := range c {
		switch b := b.(type) {
		case string:
			configString := fmt.Sprintf("%s=%s\n", a, b)
			fmt.Printf("writing config: %s\n", configString)
			file.Write([]byte(configString))
		case int:
			configString := fmt.Sprintf("%s=%d\n", a, b)
			fmt.Printf("writing config: %s\n", configString)
			file.Write([]byte(configString))
		case bool:
			configString := fmt.Sprintf("%s=%t\n", a, b)
			fmt.Printf("writing config: %s\n", configString)
			file.Write([]byte(configString))
		}
	}
	return nil
}

// genEnvVars gets all the Zookeeper environmental variables and returns them.
func genEnvVars() []string {
	limit := 99
	var envVars []string
	for ii := 1; ii <= limit; ii++ {
		envVar := baseEnv + fmt.Sprintf("%02d", ii) + baseEnvSuffix
		envVars = append(envVars, envVar)
	}
	return envVars
}

// genMyID takes the dataDir and creates the required myid file that allows Zookeeper to function.
func genMyID(dataDir string) {
	var myID string
	myID = os.Getenv("MYID")

	// if the previous file exists, then it deletes it.
	if _, err := os.Stat(fmt.Sprintf("%s", dataDir)); err == nil {
		os.RemoveAll(fmt.Sprintf("%s", dataDir))
	}

	err := os.Mkdir(fmt.Sprintf("%s", dataDir), 0777)
	if err != nil {
		fmt.Printf("error creating dataDir: %s", err)
	}
	fmt.Printf("writing myid file to %s/\n", fmt.Sprintf("%s", dataDir))

	idFile := fmt.Sprintf("%s/myid", dataDir)
	file, err := os.OpenFile(idFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	defer file.Close()

	if err != nil {
		fmt.Printf("error opening file: %s\n", err)
		fmt.Println("no file written.")
	}
	_, err = file.Write([]byte(myID))
	if err != nil {
		fmt.Printf("error writing myID: %s\n", err)
	}
}
