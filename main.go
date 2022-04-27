package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	endpoints "example.com/endpoints"
	"example.com/excelporter"
	mystructs "example.com/mysctructs"

	"github.com/darahayes/go-boom"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/xid"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

/* TODO
	FIX setID
/*


/*
BASE_URL=http://localhost
BASE_PORT=10000
REACT_URL=http://localhost:3000
*/
var BASE_URL string
var BASE_PORT string
var REACT_URL string
var FRONT_URL string

type User = mystructs.User
type Credentials = mystructs.Credentials // uppercase json hmmhmm
type UserDetails = mystructs.UserDetails
type Session = mystructs.Session
type Movie = mystructs.Movie // Need different struct to handle requests. Maybe not...
type EditMovie = mystructs.EditMovie
type Watch = mystructs.Watch
type Comment = mystructs.Comment

var sessions = map[string]Session{}

//var MovieList []Movie
var MovieList = mystructs.MovieList
var Userlist []User

var validate *validator.Validate

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONT_URL"))
	//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "*")

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "user_name",
		Value:   "",
		Expires: time.Now(),
	})
}

// https://www.sohamkamani.com/golang/session-cookie-authentication/
func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// We then get the session from our session map
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// If the session is present, but has expired, we can delete the session, and return
	// an unauthorized status
	if userSession.IsExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// REFRESH COOKIE EXPIRATION
	newSessionToken := xid.New().String()
	expiresAt := time.Now().Add(20 * time.Second)

	sessions[newSessionToken] = Session{
		Username:   userSession.Username,
		Expiration: expiresAt,
	}
	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: expiresAt,
	})
	// END OF REFRESHING PART

	// If the session is valid, return the welcome message to the user
	w.Write([]byte(fmt.Sprintf("Welcome %s!", userSession.Username)))
}

func Signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONT_URL"))
	//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "*")

	if r.Method == "POST" {
		validate = validator.New()

		creds := &Credentials{}
		err := json.NewDecoder(r.Body).Decode(creds)
		if err != nil {
			boom.BadRequest(w, "Credentials were bad")
			return
		}
		err = validate.Struct(creds)
		if err != nil {
			boom.BadRequest(w, "Credentials were bad")
			return
		}
		suser, err := getUserFromList(Userlist, creds.Username)
		if err != nil {
			boom.BadRequest(w, "No User/Password found.")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(suser.PasswordHash), []byte(creds.Password))
		if err != nil {
			boom.BadRequest(w, "Bad password.")
			return
		}

		sessionToken := xid.New().String()
		expiresAt := time.Now().Add(20 * time.Second)

		// Set the token in the session map, along with the session information
		sessions[sessionToken] = Session{
			Username:   creds.Username,
			Expiration: expiresAt,
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "user_name",
			Value:   suser.Username,
			Expires: expiresAt,
		})
		type j struct {
			User_name string `json:"user_name"`
		}
		jdata := j{User_name: suser.Username}
		json.NewEncoder(w).Encode(jdata)
	}
}

func getUserFromList(ul []User, username string) (*User, error) {
	for _, u := range ul {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, errors.New("User not found.")
}

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONT_URL"))
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "*")

	if r.Method == "POST" {
		validate = validator.New()

		creds := &Credentials{}
		err := json.NewDecoder(r.Body).Decode(creds)
		if err != nil {
			boom.BadRequest(w, "Credentials were bad")
			return
		}
		err = validate.Struct(creds)
		if err != nil {
			boom.BadRequest(w, "Credentials were bad")
			return
		}
		newUser := User{
			Username: creds.Username,
		}
		passHash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
		newUser.PasswordHash = string(passHash)
		fmt.Println("new user:", newUser)
		err = checkIfUserExists(Userlist, newUser)
		if err != nil {
			fmt.Println("error", err)
			boom.BadRequest(w, "User already exists.")
			return
		}
		newUser.SetId()
		Userlist = append(Userlist, newUser)

		fmt.Println("Users in LIST:")
		for _, u := range Userlist {
			fmt.Printf("   %s\n", u.Username)
			fmt.Printf("   %s\n", u.PasswordHash)
		}

		w.WriteHeader(http.StatusOK)
	}
}

/*
func (user *User) SetId() {
	user.Id = xid.New()
}
*/
func checkIfUserExists(userlist []User, user User) error {
	for _, u := range userlist {
		if u.Username == user.Username {
			return errors.New("User already exists.")
		}
	}
	return nil
}

func getMoviesFromFile() {
	fmt.Println("Getting movies from file")

	f, err := excelize.OpenFile("Medialists.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	/*
		cell, err := f.GetCellValue("Movies", "A4")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(cell)
	*/
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Movies")
	if err != nil {
		fmt.Println(err)
		return
	}
	for is, row := range rows {
		if is > 3 {
			//					for ir, colCell := range row {
			// }
			var newMovie Movie
			newMovie.SetId()
			newMovie.Name = row[0]
			newMovie.Year = 0
			rating, e := strconv.ParseFloat(row[1], 10)
			if e != nil {
				fmt.Println("e", e)
			}
			newMovie.Rating = int(rating)
			if len(row) == 5 {
				//fmt.Printf("is?\n %s,\n %s,\n %s,\n %s,\n %s,\n", row[0], row[1], row[2], row[3], row[4])

				if len(row[2]) > 0 {
					//fmt.Println(row[2])
					newMovie.Review = row[2]
				}
				if len(row[3]) > 0 {
					//fmt.Printf("\n\n-%s-\n", row[3])
					var newWatch Watch
					wDate, _ := time.Parse("2006-01-02", strings.Replace(row[3], ".", "-", -1))
					//fmt.Println("wDate", wDate)
					newWatch.Date = wDate
					newWatch.SetId()
					//fmt.Println("enwWathc", newWatch)
					newMovie.Watches = append(newMovie.Watches, newWatch)
				}
				if len(row[4]) > 0 {
					//fmt.Printf("\n\n-%s-\n", strings.Replace(row[4], ".", "-", -1))
					var newWatch Watch
					wDate, _ := time.Parse("2006-01-02", strings.Replace(row[4], ".", "-", -1))
					//fmt.Println("wDate", wDate)
					newWatch.Date = wDate
					newWatch.SetId()
					//fmt.Println("enwWathc", newWatch)
					newMovie.Watches = append(newMovie.Watches, newWatch)
				}
				var newWatch Watch
				newWatch.SetId()
				newMovie.Watches = append(newMovie.Watches, newWatch)
			} else if len(row) == 4 {
				//fmt.Printf("is?\n %s,\n %s,\n %s,\n %s,\n", row[0], row[1], row[2], row[3])
				if len(row[2]) > 0 {
					//fmt.Println(row[2])
					newMovie.Review = row[2]
				}
				if len(row[3]) > 0 {
					//fmt.Println(row[3])
					var newWatch Watch
					wDate, _ := time.Parse("2006-01-02", strings.Replace(row[3], ".", "-", -1))
					//fmt.Println("wDate", wDate)
					newWatch.Date = wDate
					newWatch.SetId()
					//fmt.Println("enwWathc", newWatch)
					newMovie.Watches = append(newMovie.Watches, newWatch)
				}
			} else if len(row) == 3 {
				//fmt.Printf("is?\n %s,\n %s,\n %s,\n", row[0], row[1], row[2])
				if len(row[2]) > 0 {
					//fmt.Println(row[2])
					newMovie.Review = row[2]
				}
				var newWatch Watch
				newWatch.SetId()
				newMovie.Watches = append(newMovie.Watches, newWatch)
			} else {
				//fmt.Printf("is?\n %s,\n %s,\n", row[0], row[1])
				var newWatch Watch
				newWatch.SetId()
				newMovie.Watches = append(newMovie.Watches, newWatch)
			}

			//fmt.Println("newMovie", newMovie)
			MovieList = append(MovieList, newMovie)
			//fmt.Println()
		}
	}

	fmt.Println("Got movies from file")
}

func toLowerCase(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home")
}

func movies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Movies sent")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MovieList)
}

// ENDPOINT Lisää listan
func addMovie(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

	if r.Method == "POST" {

		validate = validator.New()

		reqBody, _ := ioutil.ReadAll(r.Body)
		fmt.Println("reqBody ", string(reqBody))
		var newMovie Movie
		json.Unmarshal(reqBody, &newMovie)
		err := validate.Struct(newMovie)
		if err != nil {
			boom.BadRequest(w, "Movie not validated")
			return
		}
		fmt.Printf("new Movie:\n{\n Name: %s,\n Year: %d\n Rating: %d\n Review: %s\n}\n", newMovie.Name, newMovie.Year, newMovie.Rating, newMovie.Review)
		newMovie.SetId()
		var newWatch Watch // NOT ADDING "CANT REMEMEBER" TÄPPÄ
		json.Unmarshal(reqBody, &newWatch)
		if !newWatch.Date.IsZero() {
			newWatch.SetId()
			newMovie.Watches = append(newMovie.Watches, newWatch)
		}
		fmt.Printf("new Movie details: \nName: %s \nReview: %s\nRating: %d\nDate: %s\nPlace: %s\nNote: %s\n ", newMovie.Name, newMovie.Review, newMovie.Rating, newWatch.Date, newWatch.Place, newWatch.Note)

		/*
			err := verifyCoffee(newMvo)
			if err != nil {
				boom.BadData(w, err)
			} else {

				err = writeNewCoffee(w, newCoffee)
				if err != nil {
					boom.Internal(w, "Error while trying to create new Coffee")
				}
			}
		*/

		MovieList = append(MovieList, newMovie)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newMovie)
	}
}

func addViewing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

	if r.Method == "POST" {

		validate = validator.New()

		vars := mux.Vars(r)
		id := vars["id"]

		if !containsMovieById(MovieList, id) {
			boom.BadRequest(w, "No Movie by that ID")
		} else {
			reqBody, _ := ioutil.ReadAll(r.Body)
			fmt.Println("reqBody ", string(reqBody))

			var newWatch Watch
			json.Unmarshal(reqBody, &newWatch)
			newWatch.SetId()

			/* Tää ei toimi näin Time.time ei toimi validoiniissa näin
			err := validate.Struct(newWatch)
			if err != nil {
				boom.BadRequest(w, "Watch details not validated")
			}
			*/
			mxidm, _ := xid.FromString(id)
			movieIndex, _ := getMovieIndexFromList(mxidm)

			fmt.Println("Movielist index movie: ", MovieList[movieIndex].Name)
			MovieList[movieIndex].Watches = append(MovieList[movieIndex].Watches, newWatch)

			fmt.Printf("new Viewing details: \nMovie Name: %s \nDate: %s \nPlace: %s\nNote: %s\n", MovieList[movieIndex].Name, newWatch.Date, newWatch.Place, newWatch.Note)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newWatch)

		}
	}
}

/*
func (movie *Movie) SetId() {
	movie.Id = xid.New()
}
*/
func getMovieIndexFromList(xid xid.ID) (int, error) {
	for i, m := range MovieList {
		if xid == m.Id {
			return i, nil
		}
	}
	return -1, errors.New("No movie found by that Id")
}

func containsMovieById(ml []Movie, id string) bool {
	for _, m := range ml {
		if (m.Id).String() == id {
			return true
		}
	}
	return false
}

/*
func (watch *Watch) SetId() {
	watch.Id = xid.New()
}
*/
func getWatchIndexFromList(wl []Watch, xid xid.ID) (int, error) {
	for i, w := range wl {
		if xid == w.Id {
			return i, nil
		}
	}
	return -1, errors.New("No watch found by that Id")
}

func containsWatchById(wl []Watch, id string) bool {
	for _, w := range wl {
		if (w.Id).String() == id {
			return true
		}
	}
	return false
}

func logRequest(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request:\n       Method: %s\n       URL: %s\n       Referer: %s\n", r.Method, r.URL, r.Referer())
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func editMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT")

	if r.Method == "PUT" {

		validate = validator.New()

		vars := mux.Vars(r)
		id := vars["id"]
		fmt.Printf("PUT ID: %s\n", id)

		if !containsMovieById(MovieList, id) {
			boom.BadRequest(w, "No Movie by that ID")
		} else {
			reqBody, _ := ioutil.ReadAll(r.Body)
			fmt.Println("reqBody ", string(reqBody))

			var newMovie EditMovie
			json.Unmarshal(reqBody, &newMovie)
			err := validate.Struct(newMovie)
			if err != nil {
				boom.BadRequest(w, "Editing details not validated")
				return
			}
			mxidm, _ := xid.FromString(id)
			movieIndex, _ := getMovieIndexFromList(mxidm)

			fmt.Println("Movielist index movie: ", MovieList[movieIndex].Name)

			MovieList[movieIndex].ModifyMovie(newMovie)
			// DO SOMETHING HERE

			fmt.Printf("new Edited details: \n   Movie Name: %s \n   Year: %d\n   Rating: %d\n   Review: %s\n ", MovieList[movieIndex].Name, MovieList[movieIndex].Year, MovieList[movieIndex].Rating, MovieList[movieIndex].Review)

			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(MovieList[movieIndex])

		}
	}
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("GET ID: %s\n", id)

	if !containsMovieById(MovieList, id) {
		boom.BadRequest(w, "No Movie by that ID")
	} else {
		reqBody, _ := ioutil.ReadAll(r.Body)
		fmt.Println("reqBody ", string(reqBody))

		mxidm, _ := xid.FromString(id)
		movieIndex, _ := getMovieIndexFromList(mxidm)

		fmt.Println("Movielist index movie: ", MovieList[movieIndex].Name)

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MovieList[movieIndex])

	}
}

/*
func (m *Movie) modifyMovie(newMovie EditMovie) {
	fmt.Printf("new name: %s, new year: %d\n", (&newMovie).Name, &newMovie.Year)
	m.Name = newMovie.Name
	m.Year = newMovie.Year
	m.Rating = newMovie.Rating
	m.Review = newMovie.Review
}
*/ // NEED TO HANDLE THESE POINTER FUNCTIONS PROPERLY
func removeWatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, DELETE")

	if r.Method == "DELETE" {

		vars := mux.Vars(r)
		vid := vars["vid"]
		mid := vars["mid"]
		fmt.Printf("DELETE ID: %s\n", vid)

		if !containsMovieById(MovieList, mid) {
			boom.BadRequest(w, "No Movie by that ID")
		} else {

			mxidm, _ := xid.FromString(mid)
			movieIndex, _ := getMovieIndexFromList(mxidm)

			if !containsWatchById(MovieList[movieIndex].Watches, vid) {
				boom.BadRequest(w, "No Viewing by that ID")
			} else {
				// reqBody, _ := ioutil.ReadAll(r.Body)
				// fmt.Println("reqBody ", string(reqBody))

				vxidm, _ := xid.FromString(vid)
				watchIndex, _ := getWatchIndexFromList(MovieList[movieIndex].Watches, vxidm)

				// fmt.Println("Movielist index movie: ", watchList[watchIndex].Name)
				deletedWatch := MovieList[movieIndex].Watches[watchIndex]
				MovieList[movieIndex].Watches = append(MovieList[movieIndex].Watches[:watchIndex], MovieList[movieIndex].Watches[watchIndex+1:]...)

				w.Header().Add("Content-Type", "application/json")
				json.NewEncoder(w).Encode(deletedWatch)

			}
		}
	}
}

/*
	-Watch tietojen editointi
	-Tarkista, että editoidessa, muuta tiedot on samoja, muuten epäonnistu.
	-

*/

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	//router.HandleFunc("/test", movieservice.Test)
	router.HandleFunc("/movies", movies)
	//router.HandleFunc("/movies", endpoints.Movies) // Not working. Needs to get MOVIE INFO FROM SOMEWHERE
	router.HandleFunc("/nothing", endpoints.Nothing)
	router.HandleFunc("/movies/add", addMovie).Methods("POST", "OPTIONS")
	router.HandleFunc("/movies/{id}/viewing/add", addViewing).Methods("POST", "OPTIONS")
	router.HandleFunc("/movies/{id}/edit", editMovie).Methods("PUT", "OPTIONS")
	router.HandleFunc("/movies/{id}", getMovieById).Methods("GET", "OPTIONS")
	router.HandleFunc("/movies/{mid}/viewing/{vid}/delete", removeWatch).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/signup", Signup).Methods("POST", "OPTIONS")
	router.HandleFunc("/signin", Signin).Methods("POST", "OPTIONS")
	router.HandleFunc("/w", Welcome).Methods("GET", "OPTIONS")
	router.HandleFunc("/logout", Logout).Methods("GET", "OPTIONS")

	/*
		originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
		headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin: *"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	*/
	// log.Fatal(http.ListenAndServe(":10000", toLowerCase(handlers.CORS(originsOk, headersOk, methodsOk)(router))))
	log.Fatal(http.ListenAndServe(":10000", toLowerCase(logRequest(router))))
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	BASE_URL = os.Getenv("BASE_URL")
	BASE_PORT = os.Getenv("BASE_PORT")
	REACT_URL = os.Getenv("REACT_URL")
	FRONT_URL = os.Getenv("FRONT_URL")

	fmt.Println("base_url", BASE_URL)
	fmt.Println("base_port", BASE_PORT)
	fmt.Println("react_url", REACT_URL)
	fmt.Println("front_url", FRONT_URL)

	fmt.Println("Setting up a server on with port :10000")
	excelporter.Excelimporter()
	excelporter.Main()

	getMoviesFromFile()
	handleRequests()
}
