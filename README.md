# 16 Projects in 16 Weeks #5: Websocket Chat

## Project Overview

This project is a websocket-based chat server which people can connect to through the terminal and chat with others. 

## Learning Goals

My goals for this project are:

1. Be able to develop something in Go without the use of an LLM
    - If I get stuck for more than an hour trying with docs and experimentation, then I will chat with the LLM, but I have configured the agent to (hopefully) not give me the answers but operate in a more socratic way
    - LLMs are killing your gains
2. Create a cool little application that someone could self-host or run in a docker container and chat with their friends. 
3. Learn more about the go standard library and idiomatic go way of doing things. 
4. Context. I want to learn the context package well and understand why people use it
5. Tests. Man I gotta start writing those. Or at least know how to use the testing package in go

## How to run

1. clone the project to your machine
2. type the command `go run ./cmd start` to start the server (this will start on localhost:8080) (I'm not changing that any time soon)
3. open another terminal window (tmux btw)
4. in that terminal run `go run ./cmd tui` or `go run ./cmd repl` (if you don't want the beautiful tui experience)

## Commands within the chat

### Create a new room
`/create <room_name>`

### Join an existing room 
`/join <room_name>`

### Switch current room in tui
`/switch <room_name>`

### Change username
`/changeUsername <new_username>`


## Known issues
There are a lot of known issues with this project. For one, I don't do any validation of commands or anything so if the server recieves something that it doesn't expect it will likely crash. Same thing for the TUI, if it gets any commands that it doesn't recognize it will either crash or send those as a chat message. 

## Future improvements that I may never do

1. Make sure that the server knows how to unregister clients when they leave. 
2. Implement room deletion and leaving
3. Print announcement messages to the screen in the tui 
