package internal

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mongodb-project/db"
	"mongodb-project/models"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection("app", "users")

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, user)
}

func BulkCreateUsers(c *gin.Context) {
	var users []models.User

	// Bind JSON array
	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty users array"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection("user", "name")

	// Convert []User → []interface{}
	docs := make([]interface{}, len(users))
	for i, user := range users {
		docs[i] = user
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"inserted_count": len(result.InsertedIDs),
		"ids":            result.InsertedIDs,
	})
}

func GetUserByID(c *gin.Context) {
	idParam := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection("user", "name")

	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUsers(c *gin.Context) {

	status := c.Query("status")
	greaterThan := c.Query("greater-than")
	fmt.Println(status)
	fmt.Println(greaterThan)

	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Params missing!"})
		return
	} else if status != "claim" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect param!"})
		return
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var users []models.User

		cursor, err := db.GetCollection("user", "name").Find(ctx, bson.M{"age": bson.M{"$gt": cast.ToInt(greaterThan)}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users!"})
			return
		}
		defer cursor.Close(ctx)

		err = cursor.All(ctx, &users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract data!"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func UpdateUserById(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	update := bson.M{"$set": bson.M{"name": user.Name, "surname": user.Surname, "age": user.Age}}

	collection := db.GetCollection("user", "name")
	updatedUser, err := collection.UpdateByID(ctx, objectID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user!"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func UpdateUsers(c *gin.Context) {

	var body struct {
		AgeCategory string `json:"age_category" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}

	if ageLt := c.Query("age_gt"); ageLt != "" {
		age, err := strconv.Atoi(ageLt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age_gt"})
			return
		}
		filter["age"] = bson.M{"$gt": age}
	}

	update := bson.M{
		"$set": bson.M{
			"age_category": body.AgeCategory,
		},
	}

	collection := db.GetCollection("user", "name")

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update users"})
		return
	}

	var updatedUsers []models.User

	cursor, err := collection.Find(ctx, bson.M{
		"age_category": body.AgeCategory,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &updatedUsers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matched":  result.MatchedCount,
		"modified": result.ModifiedCount,
		"users":    updatedUsers,
	})
}

func DeleteUserByID(c *gin.Context) {
	idParam := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection("user", "name")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user!"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// StoreUserData stores user data from form input to the database
func StoreUserData(c *gin.Context) {
	var user models.User

	// Bind JSON data from request body
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user data
	if user.Name == "" || user.Surname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and surname are required"})
		return
	}

	if user.Age < 0 || user.Age > 150 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Age must be between 0 and 150"})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the collection
	collection := db.GetCollection("app", "users")

	// Insert user into database
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user in database"})
		fmt.Println("Database error:", err)
		return
	}

	// Set the ID from the insert result
	user.ID = result.InsertedID.(primitive.ObjectID)

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"id":      user.ID,
		"user":    user,
	})
}

// GetUsersByAge filters and returns users by age range
func GetUsersByAge(c *gin.Context) {
	minAge := c.Query("min_age")
	maxAge := c.Query("max_age")

	// Validate parameters
	if minAge == "" || maxAge == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_age and max_age parameters are required"})
		return
	}

	minAgeInt, err := strconv.Atoi(minAge)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_age must be a valid number"})
		return
	}

	maxAgeInt, err := strconv.Atoi(maxAge)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "max_age must be a valid number"})
		return
	}

	if minAgeInt < 0 || maxAgeInt > 150 || minAgeInt > maxAgeInt {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age range"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection("app", "users")

	// Build filter for age range
	filter := bson.M{
		"age": bson.M{
			"$gte": minAgeInt,
			"$lte": maxAgeInt,
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No users found in this age range",
			"users":   []models.User{},
			"count":   0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"users":   users,
		"count":   len(users),
	})
}
