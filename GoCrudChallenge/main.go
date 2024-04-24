package main

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "net/http"
)

type Person struct {
    ID      string   `json:"id"`
    Name    string   `json:"name" binding:"required"`
    Age     int      `json:"age" binding:"required"`
    Hobbies []string `json:"hobbies" binding:"required"`
}

var persons = []Person{
    {ID: "1", Name: "Sam", Age: 26, Hobbies: []string{}},
}

func main() {
    router := gin.Default()
    router.Use(corsMiddleware())

    
    router.GET("/person", getAllPersons)
    router.GET("/person/:personId", getPerson)
    router.POST("/person", createPerson)
    router.PUT("/person/:personId", updatePerson)
    router.DELETE("/person/:personId", deletePerson)

   
    router.NoRoute(func(c *gin.Context) {
        c.JSON(http.StatusNotFound, gin.H{"message": "Endpoint not found"})
    })

    router.Run(":3000") 
}

func getAllPersons(c *gin.Context) {
    defer func() {
        if err := recover(); err != nil {
           
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

 
    c.JSON(http.StatusOK, persons)
}

func getPerson(c *gin.Context) {
    defer func() {
        if err := recover(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    personId := c.Param("personId")
    for _, person := range persons {
        if person.ID == personId {
            c.JSON(http.StatusOK, person)
            return
        }
    }
 
    c.JSON(http.StatusNotFound, gin.H{"message": "Person not found"})
}



func createPerson(c *gin.Context) {
    defer func() {
        if err := recover(); err != nil {
        
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    var newPerson Person
    if err := c.ShouldBindJSON(&newPerson); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    newPerson.ID = uuid.NewString()
    persons = append(persons, newPerson)
    c.JSON(http.StatusOK, newPerson)
}

func updatePerson(c *gin.Context) {
    defer func() {
        if err := recover(); err != nil {
           
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    personId := c.Param("personId")
    var updatedData Person
    if err := c.ShouldBindJSON(&updatedData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    for i, person := range persons {
        if person.ID == personId {
            persons[i] = Person{ID: person.ID, Name: updatedData.Name, Age: updatedData.Age, Hobbies: updatedData.Hobbies}
            c.JSON(http.StatusOK, persons[i])
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"message": "Person not found"})
}

func deletePerson(c *gin.Context) {
    defer func() {
        if err := recover(); err != nil {
           
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    personId := c.Param("personId")
    for i, person := range persons {
        if person.ID == personId {
            persons = append(persons[:i], persons[i+1:]...)
            c.JSON(http.StatusOK, gin.H{"message": "Person deleted successfully"})
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"message": "Person not found"})
}


func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    }
}
