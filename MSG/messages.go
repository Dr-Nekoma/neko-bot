package MSG

import (
	"fmt"
	"math"
	"time"

	discord "github.com/bwmarrin/discordgo"
)

type Message struct {
	Body      string `default:""`
	Kind      string
	SubKind   string    `default:""`
	IdeaID    uint      `default:0`
	TitleLink string    `default:""`
	CreatedAt time.Time `default:""`
	Author    string    `default:""`
}

const (
	NoErrorCode       int    = 0
	ErrorCode         int    = -1
	Error             string = "error"
	Jobs              string = "jobs"
	LackOfJobs        string = "lackOfJobs"
	Project           string = "project"
	ProjectList       string = "list"
	ProjectEmptyList  string = "emptyList"
	ProjectAdd        string = "add"
	ProjectDeleteIdea string = "deleteIdea"
	ProjectDeleteId   string = "deleteId"
	Help              string = "help"
	Authorship        string = "Made by Dr.Nekoma"
)

func LackOfJobsMessage() *discord.MessageEmbed {
	return &discord.MessageEmbed{
		Title: "Try again next time!",
		Type:  "article",
		Color: 0x00acd7,
		Footer: &discord.MessageEmbedFooter{
			Text: Authorship,
		},
		Description: "I couldn't find any jobs with this key sentence!",
	}
}

func projectIdeaHelpMessage() string {
	msg := "neko!project list\n"
	msg += "neko!project add String\n"
	msg += "neko!project deleteId ID\n"
	msg += "neko!project deleteIdea String\n"
	return msg
}

func HelpMessage(descr string) *discord.MessageEmbed {
	return &discord.MessageEmbed{
		Title: "Help | Commands",
		Type:  "article",
		Color: 0x00acd7,
		Footer: &discord.MessageEmbedFooter{
			Text: Authorship,
		},
		Description: "Hey, let me help you! Here are my commands:",
		Fields: []*discord.MessageEmbedField{
			{
				Name:  "Jobs",
				Value: "neko!jobs N String",
			},
			{
				Name:  "Project Idea",
				Value: projectIdeaHelpMessage(),
			},
			{
				Name:  "Help",
				Value: "neko!help",
			},
		},
	}
}

func ErrorMessage(descr string) *discord.MessageEmbed {
	return &discord.MessageEmbed{
		Title: "Error!",
		Type:  "article",
		Color: 0xff0000,
		Footer: &discord.MessageEmbedFooter{
			Text: Authorship,
		},
		Description: descr,
	}
}

func ProjectAddMessage(msg Message, author string) *discord.MessageEmbed {
	return &discord.MessageEmbed{
		Title: "Project Idea Creation",
		Type:  "article",
		Color: 0x00acd7,
		Footer: &discord.MessageEmbedFooter{
			Text: Authorship,
		},
		Description: msg.Body,
		Fields: []*discord.MessageEmbedField{
			{
				Name:  "Author",
				Value: author,
			},
			{
				Name:  "Idea ID",
				Value: fmt.Sprint(msg.IdeaID),
			},
		},
	}
}

func ProjectEmptyListMessage() *discord.MessageEmbed {
	return &discord.MessageEmbed{
		Title: "Empty in Ideas!",
		Type:  "article",
		Color: 0x00acd7,
		Footer: &discord.MessageEmbedFooter{
			Text: Authorship,
		},
		Description: "There aren't any ideas saved!",
	}
}

func ProjectListMessage(msg Message) *discord.MessageEmbed {
	return &discord.MessageEmbed{
		Title: "Project Idea",
		Type:  "article",
		Color: 0x00acd7,
		Footer: &discord.MessageEmbedFooter{
			Text: Authorship,
		},
		Description: msg.Body,
		Fields: []*discord.MessageEmbedField{
			{
				Name:  "Author",
				Value: msg.Author,
			},
			{
				Name:  "Idea ID",
				Value: fmt.Sprint(msg.IdeaID),
			},
			{

				Name:  "Created At",
				Value: msg.CreatedAt.Format("2006-01-02 15:04"),
			},
		},
	}
}

func ProjectDeleteIdMessage(msg Message, author string) *discord.MessageEmbed {
	return &discord.MessageEmbed{
		Title: "Project Idea Deleted by ID",
		Type:  "article",
		Color: 0x00acd7,
		Footer: &discord.MessageEmbedFooter{
			Text: Authorship,
		},
		Description: msg.Body,
		Fields: []*discord.MessageEmbedField{
			{
				Name:  "Idea ID",
				Value: fmt.Sprint(msg.IdeaID),
			},
			{
				Name:  "Deleted by",
				Value: author,
			},
		},
	}
}

func ProjectDeleteIdeaMessage(msg Message, author string) *discord.MessageEmbed {
	return &discord.MessageEmbed{
		Title: "Project Idea Deleted",
		Type:  "article",
		Color: 0x00acd7,
		Footer: &discord.MessageEmbedFooter{
			Text: Authorship,
		},
		Description: msg.Body,
		Fields: []*discord.MessageEmbedField{
			{
				Name:  "Deleted by",
				Value: author,
			},
		},
	}
}

func JobMessage(titleLink string, descr string) *discord.MessageEmbed {
	return &discord.MessageEmbed{
		URL:   titleLink,
		Title: "Job Hunting",
		Type:  "article",
		Color: 0x00acd7,
		Footer: &discord.MessageEmbedFooter{
			Text: Authorship,
		},
		Image: &discord.MessageEmbedImage{
			URL: "https://jayclouse.com/wp-content/uploads/2019/06/hacker_news-1000x525-1.jpg",
		},
		Description: descr[:(uint(math.Min(float64(len(descr)), float64(4096))))],
	}
}
