package BOT

import (
	"errors"
	"log"
	"neko-bot/DB"
	"neko-bot/MSG"
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

func ParseCommand(str string, username string) []MSG.Message {
	commandArray := strings.SplitN(str, " ", 2)

	code, err := parseCommandKind(commandArray[0])

	if err != nil || len(commandArray) == 1 {
		log.Print(err)
		if code != MSG.NoErrorCode {
			return []MSG.Message{{Kind: MSG.Help}}
		}
	} else {
		if code != MSG.NoErrorCode {
			return executeCommand(code, commandArray[1], username)
		} else {
			return []MSG.Message{}
		}
	}
	return []MSG.Message{}
}

func executeCommand(code int, args string, username string) []MSG.Message {
	switch code {
	case jobSearch().code:
		jobArgs := strings.SplitN(args, " ", 2)
		howMany, err := strconv.Atoi(jobArgs[0])
		if err != nil || howMany <= 0 {
			log.Print(err)
			return []MSG.Message{{Body: "Only small positive numbers are valid!", Kind: MSG.Error}}
		}
		if len(jobArgs) < 2 {
			log.Print(err)
			return []MSG.Message{{Body: "You forgot the key sentence!", Kind: MSG.Error}}
		}
		return HackerNewsJobs(jobArgs[1], howMany)
	case storeProjectIdea().code:
		projectArgs := strings.SplitN(args, " ", 2)
		switch projectArgs[0] {
		case MSG.ProjectAdd:
			return DB.CreateIdea(projectArgs, username)
		case MSG.ProjectDeleteId:
			return DB.DeleteById(projectArgs)
		case MSG.ProjectDeleteIdea:
			return DB.DeleteByIdea(projectArgs)
		case MSG.ProjectList:
			return DB.GetIdeas()
		default:
			return []MSG.Message{{Body: "Invalid operation for project!", Kind: MSG.Error}}
		}
	case help().code:
		return []MSG.Message{{Kind: MSG.Help}}
	}
	return []MSG.Message{}
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
			return MSG.ErrorCode, errors.New("Invalid command name!")
		} else {
			return MSG.NoErrorCode, nil
		}

	}
}
