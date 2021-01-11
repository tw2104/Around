package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "path/filepath"
    "regexp"
    "time"

    jwt "github.com/form3tech-oss/jwt-go"
)

var (
    mediaTypes = map[string]string{
        ".jpeg": "image",
        ".jpg": "image",
        ".gif": "image",
        ".png": "image",
        ".mov": "video",
        ".mp4": "video",
        ".avi": "video",
        ".flv": "video",
        ".wmv": "video",
    }
)

var mySigningKey = []byte("secret")

func signupHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signup request")
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

    if r.Method == "OPTIONS" {
        return
    }

    decoder := json.NewDecoder(r.Body)
    var user User
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return 
    }

    if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        fmt.Printf("Invalid username or password\n")
        return
    }

    success, err := addUser(&user)
    if err != nil {
        http.Error(w, "Failed to save user to ElasticSearch", http.StatusInternalServerError)
        fmt.Printf("Failed to save user to ElasticSearch %v\n", err)
        return
    }
    
    if !success {
        http.Error(w, "User already exists", http.StatusBadRequest)
        fmt.Printf("User already exists\n")
        return
    }
    fmt.Printf("User added successfully: %s\n", user.Username)
}


func signinHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one sigin request")
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

    if r.Method == "OPTIONS" {
        return
    }

    // get user info from client
    decoder := json.NewDecoder(r.Body)
    var user User
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return 
    }

    exists, err := checkUser(user.Username, user.Password)
    if err != nil {
        http.Error(w, "Failed to read user from ElasticSearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read user from ElasticSearch %v\n", err)
        return 
    }

    if !exists {
        http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
        fmt.Printf("User doesn't exists or wrong password\n")
        return 
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "exp": time.Now().Add(time.Hour*24).Unix(),
    })

    tokenString, err := token.SignedString(mySigningKey)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        fmt.Printf("Failed to generate token %v\n", err)
        return 
    }

    w.Write([]byte(tokenString))
}


func uploadHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one upload request")
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

    if r.Method == "OPTIONS" {
        return
    }

    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

    p := Post {
        User: username.(string),
        Message: r.FormValue("message"),
    }

    file, header, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }

    suffix := filepath.Ext(header.Filename)
    if t, ok := mediaTypes[suffix]; ok {
        p.Type = t
    } else {
        p.Type = "unknown"
    }

    err = savePost(&p, file)
    if err != nil {
        http.Error(w, "Faild to save post to GCS or ElasticSearch", http.StatusInternalServerError)
        return
    }
    
    fmt.Println("Post is saved successfully.")
}



func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for search")
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

    if r.Method == "OPTIONS" {
        return
    }

    user := r.URL.Query().Get("user")
    keywords := r.URL.Query().Get("keywords")

    var posts []Post
    var err error
    if user != "" {
        posts, err = searchPostsByUser(user)
    } else {
        posts, err = searchPostsByKeywords(keywords)
    }

    if err != nil {
        http.Error(w, "Failed to read post from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read post from Elasticsearch %v.\n", err)
        return
    }

    js, err := json.Marshal(posts)
    if err != nil {
        http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
        fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
        return
    }
    w.Write(js)
}

