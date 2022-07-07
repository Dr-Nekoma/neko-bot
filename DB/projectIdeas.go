package DB

import (
	"fmt"
	"neko-bot/MSG"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProjectIdea struct {
	gorm.Model
	Idea   string
	Author string
}

var ConnStr string

func connectToDB() *gorm.DB {
	// Connect to database
	db, err := gorm.Open(postgres.Open(ConnStr), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&ProjectIdea{})

	return db
}

func GetIdeas() []MSG.Message {
	db := connectToDB()

	var messages []MSG.Message

	var ideas []ProjectIdea
	result := db.Find(&ideas)
	if result.Error != nil {
		return []MSG.Message{{Body: "I got an error trying to fetch ideas!", Kind: MSG.Error}}
	}
	if result.RowsAffected == 0 {
		return []MSG.Message{{Kind: MSG.Project, SubKind: MSG.ProjectEmptyList}}
	}
	for _, idea := range ideas {
		messages = append(messages, MSG.Message{Body: idea.Idea, Kind: MSG.Project, SubKind: MSG.ProjectList, IdeaID: idea.ID, CreatedAt: idea.CreatedAt, Author: idea.Author})
	}

	return messages
}

func CreateIdea(args []string, author string) []MSG.Message {
	if len(args) == 1 {
		return []MSG.Message{{Body: "You forgot to add a project idea!", Kind: MSG.Error}}
	}
	idea := args[1]
	db := connectToDB()

	record := ProjectIdea{Idea: idea, Author: author}

	result := db.Create(&record)

	if result.Error != nil {
		return []MSG.Message{{Body: "I got an error trying to add an idea!", Kind: MSG.Error}}
	}
	return []MSG.Message{{Body: fmt.Sprintf("Idea '%s' by %s has been created!", idea, author), Kind: MSG.Project, SubKind: MSG.ProjectAdd}}
}

func DeleteById(args []string) []MSG.Message {
	if len(args) == 1 {
		return []MSG.Message{{Body: "You forgot to specify the ID!", Kind: MSG.Error}}
	}
	id := args[1]
	db := connectToDB()

	result := db.Unscoped().Delete(&ProjectIdea{}, id)

	if result.Error != nil {
		return []MSG.Message{{Body: "I got an error trying to delete this idea via ID!", Kind: MSG.Error}}
	}
	if result.RowsAffected == 0 {
		return []MSG.Message{{Body: "There isn't an idea with this ID to be deleted!", Kind: MSG.Error}}
	}

	return []MSG.Message{{Body: fmt.Sprintf("Idea with ID %s has been deleted!", id), Kind: MSG.Project, SubKind: MSG.ProjectDeleteId}}
}

func DeleteByIdea(args []string) []MSG.Message {
	if len(args) == 1 {
		return []MSG.Message{{Body: "You forgot to specify the idea!", Kind: MSG.Error}}
	}
	idea := args[1]
	db := connectToDB()

	result := db.Where("idea = ?", idea).Unscoped().Delete(&ProjectIdea{})

	if result.Error != nil {
		return []MSG.Message{{Body: "I got an error trying to delete this idea via its content!", Kind: MSG.Error}}
	}
	if result.RowsAffected == 0 {
		return []MSG.Message{{Body: "There isn't an idea with this name to be deleted!", Kind: MSG.Error}}
	}

	return []MSG.Message{{Body: fmt.Sprintf("Idea '%s' has been deleted!", idea), Kind: MSG.Project, SubKind: MSG.ProjectDeleteIdea}}
}
