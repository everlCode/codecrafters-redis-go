package handlers

const (
	PX = "px"
)

type ArgParser struct {

}

func New() *ArgParser {
	return &ArgParser{}
}

func (ap ArgParser) Parse(args []string, cmd Command) map[string]int {
	// for _, v := range args {
		
	// }

	return map[string]int{
		"sd": 1,
	}
}