package BOT

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

type command struct {
	name string
	code int
}

func jobSearch() command {
	return command{name: "neko!jobs", code: 1}
}

func storeProjectIdea() command {
	return command{name: "neko!project", code: 2}
}

func help() command {
	return command{name: "neko!help", code: 3}
}

const (
	NoErrorCode int = 0
	ErrorCode   int = -1
)

func printHelp() string {
	msg := "Hey, let me help you! Here are my commands:\n"
	msg += "JOBS: neko!jobs N String\n"
	msg += "PROJECT: neko!project String\n"
	return msg
}

func ParseCommand(str string) []string {
	commandArray := strings.SplitN(str, " ", 2)

	code, err := parseCommandKind(commandArray[0])

	if err != nil || len(commandArray) == 1 {
		log.Print(err)
		if code != NoErrorCode {
			msg := printHelp()
			return []string{msg}
		}
	} else {
		if code != NoErrorCode {
			return executeCommand(code, commandArray[1])
		} else {
			return []string{}
		}
	}
	return []string{}
}

func executeCommand(code int, args string) []string {
	switch code {
	case jobSearch().code:
		jobArgs := strings.SplitN(args, " ", 2)
		howMany, err := strconv.Atoi(jobArgs[0])
		if err != nil {
			log.Print(err)
			return []string{}
		} else {
			return HackerNewsJobs(jobArgs[1], howMany)
		}
	case storeProjectIdea().code:
		return []string{"TODO: store the ideas in the database"}
	case help().code:
		return []string{printHelp()}
	}
	return []string{}
}

func parseCommandKind(kind string) (int, error) {
	switch kind {
	case jobSearch().name:
		return jobSearch().code, nil
	case storeProjectIdea().name:
		return storeProjectIdea().code, nil
	case help().name:
		return help().code, nil
	default:
		if strings.Contains(kind, "neko!") {
			return ErrorCode, errors.New("Invalid command name!")
		} else {
			return NoErrorCode, nil
		}

	}
}
