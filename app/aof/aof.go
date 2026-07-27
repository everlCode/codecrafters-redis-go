package aof

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Aof struct {
	Config *config.Config
}

func New(config *config.Config) *Aof {
	if config.Appendonly == "yes" {
		createAppendOnlyDirectory(config)
	}


	return &Aof{
		Config: config,
	}
}

func incrFilePath(config *config.Config) string {
	dirName := filepath.Join(config.Dir, config.Appenddirname)
	return filepath.Join(dirName, config.Appendfilename+".1.incr.aof")
}

func createAppendOnlyDirectory(config *config.Config) {
	dirName := filepath.Join(config.Dir, config.Appenddirname)

	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = os.Create(incrFilePath(config))
	if err != nil {
		fmt.Println(err.Error())
	}

	filename := filepath.Join(dirName, config.Appendfilename+".manifest")
	manifestFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	manifestFile.WriteString(fmt.Sprintf("file %s.1.incr.aof seq 1 type i", config.Appendfilename))
}

func (aof *Aof) Write(command resp.Value) {
	file, err := os.OpenFile(incrFilePath(aof.Config), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	defer file.Close()
	file.Write(command.Marshal())
}