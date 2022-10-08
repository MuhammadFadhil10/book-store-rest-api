package controller

import (
	"encoding/json"
	"go-rest/config"
	"go-rest/models"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	
	var user models.User
	decodeErr := json.NewDecoder(r.Body).Decode(&user)


	if decodeErr != nil {
		log.Fatal(decodeErr)
	}

	resp := make(map[string]interface{})

	userExist := user.CheckExistingEmail(user.Email)

	if userExist["isExist"].(bool) {
		resp["status"] = false
		resp["message"] = "This email already used!"
		
		json.NewEncoder(w).Encode(resp)
		return
	}
	
	user.Id = uuid.New().String()
	// fmt.Println(user.Id)

	hashedPassword,hashErr := bcrypt.GenerateFromPassword([]byte(user.Password),12)

	if hashErr != nil {
		log.Fatal(hashErr)
	}



	user.Password = string(hashedPassword)
	

	insertErr := user.CreateUser(user.Id,user.Name, user.Email, user.Password)
	if insertErr != nil {
		log.Fatal(insertErr.Error())
	}

	
	resp["status"] = true
	resp["message"] = "Register success!"

	json.NewEncoder(w).Encode(resp)

}

func Login(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	
	var user models.User
	decodeErr := json.NewDecoder(r.Body).Decode(&user)

	if decodeErr != nil {
		log.Fatal(decodeErr)
	}

	resp := make(map[string]interface{})

	userExist := user.CheckExistingEmail(user.Email)

	if !userExist["isExist"].(bool) {
		resp["status"] = false
		resp["message"] = "Wrong email or password!"
		
		json.NewEncoder(w).Encode(resp)
		return
	}

	// userData := user.FindById()
	
	userData := userExist["user"].(models.User)

	passwordErr := bcrypt.CompareHashAndPassword([]byte(userData.Password),[]byte(user.Password))

	if passwordErr != nil {
		resp["status"] = false
		resp["message"] = "Wrong email or password!"
		
		json.NewEncoder(w).Encode(resp)
		return
	}

	// jwt
	// create jwt claim based on jwt claim config struct
	expTime := time.Now().Add(time.Minute * 1)
	jwtClaim := config.JWTClaim {
		UserName: userData.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// create jwt token algorithm
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
	// sign jwt
	token, err := tokenAlgo.SignedString(config.JWT_KEY)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	// set token to cookie
	http.SetCookie(w,&http.Cookie{
		Name: "access_token",
		Path: "/",
		Value: token,
		HttpOnly: true,
	})
	
	// response
	resp["status"] = true
	resp["message"] = "Login success!"
	resp["accessToken"] = token

	// json.NewEncoder(w).Encode(resp)
	response,_ := json.Marshal(resp)
	w.Write(response)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	http.SetCookie(w,&http.Cookie{
		Name: "access_token",
		Path: "/",
		Value: "",
		HttpOnly: true,
		MaxAge: -1,
	})

	resp := map[string]interface{} {
		"Status": true,
		"message": "success logout",
	}

	response,_ := json.Marshal(resp)
	w.Write(response)
}