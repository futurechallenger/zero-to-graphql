package schema

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	baseURL = ""
)

// GetJSONFromRelativeURL Get json from relative url
func GetJSONFromRelativeURL(relativeURL string) string {
	resp, err := http.Get(baseURL + relativeURL)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(body)
	return string(body[:])
}

// GetAllPeople will request the url for all people
func GetAllPeople() string {
	return GetJSONFromRelativeURL("/people/all")
}

// GetPerson will request just a person
func GetPerson(personID int64) string {
	return GetJSONFromRelativeURL("/people/" + strconv.FormatInt(personID, 10))
}

// GetFriends will request a person's all friends
func GetFriends(personID int64) string {
	return //TODO: dataloader still needed in this situation
}
