# neko-bot

This is a bot for Dr.Nekoma's discord server implemented in Go.

## External Dependencies

- [discordgo](https://github.com/bwmarrin/discordgo)
- [godotenv](https://github.com/joho/godotenv)
- [html-to-markdown](https://pkg.go.dev/github.com/JohannesKaufmann/html-to-markdown)
- [gorm](https://gorm.io/index.html)

## Neko-bot's Functionalities

This bot is capable of:

- Search up to an arbitrary number of jobs in HackerNews website using a key sentence
- Accumulate project ideas in a Heroku database (Create, Read and Delete capabilities)

Available commands:

```
neko!jobs 3 Haskell
neko!project add Make a distributed system in OCaml
neko!project list
neko!project deleteIdea Make a distributed system in OCaml
```

## Developers

- EduardoLR10
- ribeirotomas1904
- MMagueta

## Dr.Nekoma

Builded live on [twitch](https://www.twitch.tv/drnekoma) and archived on [youtube](https://www.youtube.com/channel/UCMyzdYsPiBU3xoqaOeahr6Q)
