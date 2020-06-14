package main

import (
	"math/rand"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type payloadEvent struct {
	Event string
	Asset string
}

type wsclient struct {
	upgrader websocket.Upgrader
	users    []*User
}

func (p *wsclient) handle(w http.ResponseWriter, r *http.Request) {

	//Attempt to upgrade the connection
	c, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error: ", err)
		return
	}

	//Add the connection
	user := p.addUser(c)

	//Defer the close until the end
	defer c.Close()
	defer p.removeUser(user)

	for {

		//Prepare the payload to read into
		payload := &payloadEvent{}

		//Read the payload, if one is available
		err := c.ReadJSON(payload)
		if err != nil {
			log.Println("error on read: ", err)
			return
		}
	}
}

//Sends a message to all users
func (p *wsclient) Broadcast(e payloadEvent) int {
	log.Println("Broaadcast: ", e)
	count := 0
	for _, user := range p.users {

		//Send  the payload and log any errors
		err := user.send(e)
		if err != nil {
			log.Println(user.identifier, " error on write: ", err)
		} else {
			//No errors so increment our count
			count++
		}
	}
	return count
}

//addUser adds a user to the list. Does not send events because the user is requried to have authorized
func (p *wsclient) addUser(websocket *websocket.Conn) *User {
	user := newUser(websocket)
	//log.Println("New user connected: ", user.identifier)
	p.users = append(p.users, user)
	return user
}

//removeUser will delete a user from the list.
func (p *wsclient) removeUser(user *User) {
	for i, u := range p.users {
		if u.identifier == user.identifier {
			copy(p.users[i:], p.users[i+1:])
			p.users[len(p.users)-1] = nil
			p.users = p.users[:len(p.users)-1]
			break
		}
	}

	//log.Println("User Disconnected", user.identifier)
	//Tell everyone this user disconnected
	//p.Broadcast(UserIdentifyEvent{Identity: user.identifier, DisplayName: user.displayname, State: UserIdentifyDisconnect})
}

// User represents a connection to the project
type User struct {
	identifier  string
	displayname string
	authorized  bool
	websocket   *websocket.Conn
}

func newUser(websocket *websocket.Conn) *User {
	return &User{
		authorized:  true,
		identifier:  strconv.Itoa(rand.Int()),
		displayname: "User",
		websocket:   websocket,
	}
}

func (user *User) send(e payloadEvent) error {
	if !user.authorized {
		return errors.New("not authorized")
	}

	return user.websocket.WriteJSON(e)
}

func (user *User) close() error {
	return user.websocket.Close()
}

