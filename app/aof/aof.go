package aof

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Aof struct {
	Config *config.Config
}

func New(config *config.Config) *Aof {
	if config.Appendonly == "yes" {
		if _, err := os.Stat(manifestFilePath(config)); os.IsNotExist(err) {
			createAppendOnlyDirectory(config)
		}
	}

	return &Aof{
		Config: config,
	}
}

func dirPath(config *config.Config) string {
	return filepath.Join(config.Dir, config.Appenddirname)
}

func incrFilePath(config *config.Config) string {
	return filepath.Join(dirPath(config), config.Appendfilename+".1.incr.aof")
}

func manifestFilePath(config *config.Config) string {
	return filepath.Join(dirPath(config), config.Appendfilename+".manifest")
}

func createAppendOnlyDirectory(config *config.Config) {
	err := os.MkdirAll(dirPath(config), 0755)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = os.Create(incrFilePath(config))
	if err != nil {
		fmt.Println(err.Error())
	}

	manifestFile, err := os.Create(manifestFilePath(config))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer manifestFile.Close()
	manifestFile.WriteString(fmt.Sprintf("file %s.1.incr.aof seq 1 type i", config.Appendfilename))
}

// readManifest returns the append-only file names listed in the manifest, in order.
func readManifest(config *config.Config) ([]string, error) {
	data, err := os.ReadFile(manifestFilePath(config))
	if err != nil {
		return nil, err
	}

	var files []string
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0] != "file" {
			continue
		}
		files = append(files, fields[1])
	}

	return files, nil
}

// Load reads every file referenced by the manifest and returns the commands
// stored in them, in the order they were written.
func (aof *Aof) Load() []resp.Value {
	var commands []resp.Value

	if aof.Config.Appendonly != "yes" {
		return commands
	}

	files, err := readManifest(aof.Config)
	if err != nil {
		fmt.Println(err.Error())
		return commands
	}

	for _, name := range files {
		file, err := os.Open(filepath.Join(dirPath(aof.Config), name))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		commands = append(commands, readCommands(file)...)
		file.Close()
	}

	return commands
}

func readCommands(file *os.File) []resp.Value {
	var commands []resp.Value
	parser := resp.New(file)

	for {
		command, err := readCommand(parser)
		if err != nil {
			break
		}
		commands = append(commands, command)
	}

	return commands
}

// readCommand wraps parser.Read to turn a panic on a truncated/corrupted
// trailing command (e.g. after a crash mid-write) into a plain error, so the
// caller can stop reading that file instead of crashing the whole replay.
func readCommand(parser *resp.Parser) (value resp.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("incomplete AOF command: %v", r)
		}
	}()

	return parser.Read()
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