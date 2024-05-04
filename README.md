# Oude Markt Bar Picker V2
Author: Mathias Are

Interactive web application that helps to choose the bar to go to with friends in Leuven Oude-Markt.



### Usage
The user is expected to create a new vote from the UI and share the given vote code with his/her friends. The friends can use the same code to join the vote room. During the voting process all user votes are updated dynamically in real time via websocket message passing. Additionally the vote state is stored in the Turso DB, so the vote can later be rejoined and its result observed by visiting the correct URL.

### Idea

The idea of the project was to study and showcase the capabilities of HTMX and Go Templates in more complex and interactive scenarios like using websockets and broadcasting messages to websocket channels.

### Technolgies used

- Go + Chi framework
- HTMX + websockets extention
- Turso LibSQL database
- GORM
- TailwindCSS

Deployable as a Docker container

### Running the app

In project root Run: 
```bash 
go run main.go
```


