/*
* Copyright 2016 Comcast Cable Communications Management, LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

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
		autopurgeRetainCount int
		autopurgePurgeInterval int
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
	flag.IntVar(&autopurgeRetainCount, "autopurge.snapRetainCount", 3, "most recent snapshots and the corresponding transaction logs in the dataDir and dataLogDir respectively and deletes the rest. defaults to 3.")
	flag.IntVar(&autopurgePurgeInterval, "autopurge.purgeInterval", 72, "the time interval in hours for which the purge task has to be triggered. set to a positive integer (1 and above) to enable the auto purging. defaults to 72.")

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
	configs["autopurge.snapRetainCount"] = autopurgeRetainCount
	configs["autopurge.purgeInterval"] = autopurgePurgeInterval

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
			fmt.Printf("Writing config: %s", configString)
			file.Write([]byte(configString))
		case int:
			configString := fmt.Sprintf("%s=%d\n", a, b)
			fmt.Printf("Writing config: %s", configString)
			file.Write([]byte(configString))
		case bool:
			configString := fmt.Sprintf("%s=%t\n", a, b)
			fmt.Printf("Writing config: %s", configString)
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
	fmt.Printf("Writing myid file to %s/\n", fmt.Sprintf("%s", dataDir))

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
