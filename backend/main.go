package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Blog struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
const blogFile = "blogs.json"
const userFile = "users.json"

func checkUsers(reqUser User) bool {
	file,err := os.Open(userFile)
	var users []User

	if(err != nil){
		return false
	}
	bytes, err := io.ReadAll(file)
	if(err != nil){
		return false
	}
	err = json.Unmarshal(bytes, &users)
	if(err != nil){
		return false
	}

	for _,user := range users{
		if(user.Username == reqUser.Username && user.Password == reqUser.Password){
			return true
		}
	}
	return false
}

func checkBlogs(title string) int {
	blogs, err := loadBlogs()
	if err != nil {
		fmt.Println("Error loading blogs:", err)
		return -1
	}

	found := -1
	for indx, blog := range blogs {
		if title == blog.Title {
			found = indx
		}
	}

	return found
}

func loadBlogs() ([]Blog, error) {
	var blogs []Blog
	file, err := os.Open(blogFile)
	if err != nil {
		if os.IsNotExist(err) {
			return blogs, nil
		}
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return blogs, nil
	}
	err = json.Unmarshal(bytes, &blogs)
	return blogs, err
}

func saveBlogs(blogs []Blog) error {
	bytes, err := json.MarshalIndent(blogs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(blogFile, bytes, 0644)
}

func addBlog(blog Blog) error {
	blogs, err := loadBlogs()
	if err != nil {
		return err
	}
	isPresent := checkBlogs(blog.Title)

	if isPresent != -1 {
		return fmt.Errorf("blog with title '%s' already exists", blog.Title)
	}
	blogs = append(blogs, blog)
	return saveBlogs(blogs)
}

func removeBlog(title string) error {
	blogs, err := loadBlogs()
	if err != nil {
		return err
	}
	idx := checkBlogs(title)

	if idx == -1 {
		return fmt.Errorf("blog not found")
	}
	blogs = append(blogs[:idx], blogs[idx+1:]...)
	return saveBlogs(blogs)
}

func editBlog(title string, newBlog Blog) error {
	blogs, err := loadBlogs()
	if err != nil {
		return err
	}
	idx := checkBlogs(title)

	if idx == -1 {
		return fmt.Errorf("blog not found")
	}

	blogs[idx] = newBlog
	return saveBlogs(blogs)
}

func writeJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func addBlogHandler(w http.ResponseWriter, r *http.Request) {

	if(os.Getenv("AUTH") != "true"){
		writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var blog Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := addBlog(blog)
	if err != nil {
		// writeJSONError(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Blog added successfully"})
}

func removeBlogHandler(w http.ResponseWriter, r *http.Request) {

	if(os.Getenv("AUTH") != "true"){
		writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := removeBlog(req.Title)
	if err != nil {
		// writeJSONError(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Blog removed successfully"})
}

func editBlogHandler(w http.ResponseWriter, r *http.Request) {

	if(os.Getenv("AUTH") != "true"){
		writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var blog Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := editBlog(blog.Title, blog)
	if err != nil {
		// writeJSONError(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Blog edited successfully"})
}


func blogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	blogs, err := loadBlogs()
	if err != nil {
		http.Error(w, "Failed to load blogs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user) 

	if err != nil{
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// var found bool
	found := checkUsers(user)

	if found {
		// we set the flag to be true ig 
		os.Setenv("AUTH","true")
		fmt.Println("found the username")
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged in successfully"})
		return
	}


	// auth := os.Getenv("AUTH")
	// fmt.Println(auth)
	json.NewEncoder(w).Encode(map[string]string{"message": "Log in Unsuccessful"})
}

func logoutHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}


	os.Setenv("AUTH","false")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func main() {


	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	

	// user := User{"Shyam","Shym"}
	// fmt.Println(checkUsers(user))
	// auth := os.Getenv("AUTH")
	// fmt.Println(auth)


	http.HandleFunc("/add", withCORS(addBlogHandler))
	http.HandleFunc("/remove", withCORS(removeBlogHandler))
	http.HandleFunc("/edit", withCORS(editBlogHandler))
	http.HandleFunc("/blogs", withCORS(blogsHandler))
	http.HandleFunc("/login", withCORS(loginHandler))
	http.HandleFunc("/logout", withCORS(logoutHandler))


	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
