package main 

import (
"log"
"github.com/gorilla/mux"
"github.com/gorilla/websocket"
"net/http"
"encoding/json"
"sync"
"os"
)

type Response struct{
	Message string `json:"message"`
	Status int `json:"status"`
	IsValid bool `json:"isvalid"`
}

var Users = struct{
	m map[string] User
	sync.RWMutex
}{m: make(map[string] User)}

type User struct{
	User_Name string 
	WebSocket *websocket.Conn
}

func HolaMundo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte ("Hola Mundo desde Go"))
	
}

func HolaMundoJson(w http.ResponseWriter, r *http.Request) {
	
	response := CreateResponse("Esto esta en formato Json",200,true)
	json.NewEncoder(w).Encode(response)
}

func CreateResponse(message string , status int, valid bool) Response {
	
	return Response{message,status,valid}
}

func LoadStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r,"./Front/index.html")
}

func UserExit(user_name string)bool {
	Users.RLock();
	defer Users.RUnlock()

	if _, ok :=Users.m[user_name]; ok{
		return true
	}
	return false

}

func validate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user_name := r.FormValue("user_name")
	response := Response{}

	if UserExit(user_name) {

		response.IsValid =false

		
	}else{

		response.IsValid = true

	}
	json.NewEncoder(w).Encode(response)
}

func CreateUser(user_name string, ws *websocket.Conn) User {
	
	return User{user_name,ws}
}
func AddUser(user User) {
	Users.Lock()
	defer Users.Unlock()
	Users.m[user.User_Name] = user

	
}
func RemoveUser(user_name string) {
	log.Println("El usuario se fue")
    Users.Lock()
	defer Users.Unlock()
	delete(Users.m,user_name)
	
}
func SendMessage(type_message int, message []byte) {
	Users.RLock()
	defer Users.RUnlock()

	for _, user := range Users.m{
		 if err := user.WebSocket.WriteMessage(type_message,message); err != nil{
		 	return
		 }
	}
}
func ToArryByte(value string)[]byte {
	return []byte(value)
}

func ConactMessage(user_name string,array []byte) string {
	return user_name + ":" + string(array[:])
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_name := vars["user_name"]
	ws, err := websocket.Upgrade(w,r,nil,1024,1024)
	if err != nil {
		log.Println(err)
		return
	}

	current_user := CreateUser(user_name, ws)
	AddUser(current_user)
	log.Println("Nuevo usuario agregado")

	for  {
		type_message,message,err := ws.ReadMessage()
		if err != nil {
			RemoveUser(user_name)
			return
		}
		final_message := ConactMessage(user_name,message)
		SendMessage(type_message,ToArryByte(final_message))
	}

}
 
func main() {

	cssHandle := http.FileServer(http.Dir("./Front/CSS/"))
	Js_Handle := http.FileServer(http.Dir("./Front/JS/"))


	mux := mux.NewRouter()
	mux.HandleFunc("/Hola",HolaMundo).Methods("GET")
	mux.HandleFunc("/HolaJson", HolaMundoJson).Methods("GET")
	mux.HandleFunc("/",LoadStatic).Methods("GET")
	mux.HandleFunc("/validate",validate).Methods("POST")
	mux.HandleFunc("/chat/{user_name}",WebSocket).Methods("GET")

	port := os.Getenv("PORT")

	if port=="" {
		port= "8000"
	}

	http.Handle("/",mux)
	http.Handle("/CSS/",http.StripPrefix("/CSS/", cssHandle))
	http.Handle("/JS/",http.StripPrefix("/JS/", Js_Handle))
	log.Println("El servidor se encuentra en el puerto 8000")
	log.Fatal( http.ListenAndServe(":"+port,nil))

}