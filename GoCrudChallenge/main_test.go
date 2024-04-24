package main

import (
    "bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPersons(t *testing.T) {
 
    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/person", nil)

    
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var persons []Person
    err := json.Unmarshal(w.Body.Bytes(), &persons)
    assert.NoError(t, err)

    assert.Equal(t, persons, persons) 
}

func TestGetPerson(t *testing.T) {
    
    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/person/1", nil)

    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)

    var person Person
    err := json.Unmarshal(w.Body.Bytes(), &person)
    assert.NoError(t, err)

    assert.Equal(t, person, person)
}

func TestCreatePerson(t *testing.T) {

    router := setupRouter()
    w := httptest.NewRecorder()
    reqBody := []byte(`{"name":"John","age":30,"hobbies":["Reading","Coding"]}`)
    req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)


    var createdPerson Person
    err := json.Unmarshal(w.Body.Bytes(), &createdPerson)
    assert.NoError(t, err)

    assert.NotEmpty(t, createdPerson.ID)
    assert.Equal(t, "John", createdPerson.Name)
    assert.Equal(t, 30, createdPerson.Age)
    assert.Equal(t, []string{"Reading", "Coding"}, createdPerson.Hobbies)
}
func TestPostValidationAllEmpty(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    reqBody := []byte(`{}`)
    req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostValidationNameEmpty(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    reqBody := []byte(`{"age":26,"hobbies":[]}`)
    req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostValidationAgeEmpty(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    reqBody := []byte(`{"name":"sam","hobbies":[]}`)
    req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostValidationAgeNumber(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    reqBody := []byte(`{"name":"sam","age":"bad","hobbies":[]}`)
    req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostValidationHobbiesEmpty(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    reqBody := []byte(`{"name":"sam","age":21}`)
    req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostValidationHobbiesArray(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    reqBody := []byte(`{"name":"sam","age":21,"hobbies":"fighting"}`)
    req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePerson(t *testing.T) {
 
    router := setupRouter()
    w := httptest.NewRecorder()
    reqBody := []byte(`{"name":"Updated Name","age":35,"hobbies":["Swimming","Cooking"]}`)
    req, _ := http.NewRequest("PUT", "/person/1", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")


    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var updatedPerson Person
    err := json.Unmarshal(w.Body.Bytes(), &updatedPerson)
    assert.NoError(t, err)

    assert.Equal(t, "1", updatedPerson.ID)
    assert.Equal(t, "Updated Name", updatedPerson.Name)
    assert.Equal(t, 35, updatedPerson.Age)
    assert.Equal(t, []string{"Swimming", "Cooking"}, updatedPerson.Hobbies)
}


func TestDeletePerson(t *testing.T) {

    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("DELETE", "/person/1", nil)


    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Person deleted successfully", response["message"])
}

func TestNonExistingUser(t *testing.T) {

    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/person/100", nil)

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
}


func TestNonExistingEndpoint(t *testing.T) {

    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/test/non-existing/endpoint", nil)

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
}




func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/person", getAllPersons)
	router.GET("/person/:personId", getPerson)
	router.POST("/person", createPerson)
	router.PUT("/person/:personId", updatePerson)
	router.DELETE("/person/:personId", deletePerson)
	return router
}
