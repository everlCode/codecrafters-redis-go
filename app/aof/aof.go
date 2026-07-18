package aof

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/redis-starter-go/app/config"
)

type Aof struct {

}

func New(config *config.Config) *Aof {
	if config.Appendonly == "yes" {
		createAppendOnlyDirectory(config)
	}


	return &Aof{}
}

func createAppendOnlyDirectory(config *config.Config) {
	dirName := filepath.Join(config.Dir, config.Appenddirname)

	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		fmt.Println(err.Error())
	}
	filename := filepath.Join(dirName, config.Appendfilename)
	filename = filename + ".1.incr.aof"

	_, err = os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
}