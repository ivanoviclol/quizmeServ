package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Passwd string `json:"passwd,omitempty"`
}

type Quiz struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	UserID    string     `json:"user_id,omitempty"`
	Questions []Question `json:"questions,object"`
}

type Question struct {
	ID      string   `json:"id,omitempty"`
	Text    string   `json:"text,omitempty"`
	QuizID  string   `json:"quizid,omitempty"`
	Answers []Answer `json:"answers,object"`
}

type Answer struct {
	ID         string `json:"id,omitempty"`
	Text       string `json:"text,omitempty"`
	QuestionID string `json:"questionid,omitempty"`
	Correct    bool   `json:"correct, omitempty"`
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/quizdb")
	text := "SELECT `Name` FROM `users`;"

	rows, err := db.Query(text)
	checkErr(err)
	users := []string{}
	for rows.Next() {
		var Table_Comment string
		err = rows.Scan(&Table_Comment)
		checkErr(err)
		users = append(users, Table_Comment)
	}
	for _, user := range users {
		fmt.Printf(user)
	}
	json.NewEncoder(w).Encode(users)
}
func IfUserExists(w http.ResponseWriter, r *http.Request) {
	params := r.Header.Get("email")
	fmt.Print(params + "\n")
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/quizdb")
	text := "Select exists (SELECT email FROM `users` where `email` = '" + params + "');"
	var exists bool
	rows, err := db.Query(text)
	checkErr(err)
	checkedUser := User{}
	for rows.Next() {
		err = rows.Scan(&exists)
		checkErr(err)
		if exists {
			fmt.Print("User  exist \n")
			checkedUser.ID = "1"

		} else {
			fmt.Print("User doesnt exists \n")
			checkedUser.ID = "0"
		}
	}
	json.NewEncoder(w).Encode(checkedUser)
}
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var newUser User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	checkErr(err)
	querry := "INSERT into `users` (Name, email, password) values ('" + newUser.Name + "', '" + newUser.Email + "', '" + newUser.Passwd + "');"
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/quizdb")
	_, err = db.Exec(querry)
	checkErr(err)
}

func CreateQuiz(w http.ResponseWriter, r *http.Request) {
	user_id := r.Header.Get("user_id")
	querry := "INSERT into `quiz` (Name, user_id) values ('New Awesome Quiz', '" + user_id + "');"
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/quizdb")
	_, err = db.Exec(querry)
	checkErr(err)
	json.NewEncoder(w).Encode(user_id)
}

func UpdateQuiz(w http.ResponseWriter, r *http.Request) {
	var newQuiz Quiz
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newQuiz)
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/quizdb")
	checkErr(err)
	fmt.Print(newQuiz)
	fmt.Print("\n")
	fmt.Print(newQuiz.Questions)
	fmt.Print("\n")
	fmt.Print(newQuiz.Questions[0].Answers)
	fmt.Print("\n")

	querry := "Delete from quiz where Id = '" + newQuiz.ID + "';"
	_, err = db.Exec(querry)
	checkErr(err)

	querry = "Insert into quiz (`Name`, `user_id`) values ('" + newQuiz.Name + "', '" + newQuiz.UserID + "');"
	_, err = db.Exec(querry)
	checkErr(err)

	querry = "Select Id from quiz where `Name` = '" + newQuiz.Name + "' and `user_id` = '" + newQuiz.UserID + "';"
	rows, err := db.Query(querry)
	checkErr(err)

	indexesId := []string{}
	for rows.Next() {
		var Table_Comment string
		err = rows.Scan(&Table_Comment)
		checkErr(err)
		indexesId = append(indexesId, Table_Comment)
	}

	for number, value := range newQuiz.Questions {

		querry := "Delete from question where quiz_id = '" + newQuiz.ID + "';"
		_, err = db.Exec(querry)
		checkErr(err)

		querry = "Insert into question (`question`, `quiz_id`) values ('" + value.Text + "', '" + indexesId[0] + "');"
		_, err = db.Exec(querry)
		checkErr(err)

		querry = "Select Id from question where question = '" + value.Text + "' and quiz_id = '" + indexesId[0] + "';"
		rows, err := db.Query(querry)
		checkErr(err)

		indexes := []string{}
		for rows.Next() {
			var Table_Comment string
			err = rows.Scan(&Table_Comment)
			checkErr(err)
			indexes = append(indexes, Table_Comment)
		}

		for _, answer := range newQuiz.Questions[number].Answers {

			querry := "Delete from answer where question_id = '" + newQuiz.Questions[number].ID + "';"
			_, err = db.Exec(querry)
			checkErr(err)

			querry = "Insert into answer (`answer`, `question_id`, `correct`) values ('" + answer.Text + "', '" + indexes[0] + "', " + strconv.FormatBool(answer.Correct) + ");"
			_, err = db.Exec(querry)
			checkErr(err)
		}

	}
	json.NewEncoder(w).Encode(newQuiz)

}

func DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	var newQuiz Quiz
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newQuiz)
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/quizdb")
	checkErr(err)

	querry := "Delete from quiz where Id = '" + newQuiz.ID + "';"
	_, err = db.Exec(querry)
	checkErr(err)

	querry = "Delete from question where quiz_id = '" + newQuiz.ID + "';"
	_, err = db.Exec(querry)
	checkErr(err)

	for number, _ := range newQuiz.Questions {
		querry := "Delete from answer where question_id = '" + newQuiz.Questions[number].ID + "';"
		_, err = db.Exec(querry)
		checkErr(err)
	}

	json.NewEncoder(w).Encode(newQuiz)

}

func DeletePerson(w http.ResponseWriter, r *http.Request) {}

func GetQuiz(w http.ResponseWriter, r *http.Request) {
	user_id := r.Header.Get("user_id")
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/quizdb")
	querry := "SELECT Id, Name, user_id from `quiz` where `user_id` = '" + user_id + "';"
	rows, err := db.Query(querry)
	checkErr(err)
	quizList := []Quiz{}

	for rows.Next() {
		quiz := Quiz{}
		err = rows.Scan(&quiz.ID, &quiz.Name, &quiz.UserID)
		checkErr(err)
		quizList = append(quizList, quiz)
	}

	for i := 0; i < len(quizList); i++ {
		questionList := []Question{}
		querry := "Select id, question from question where  quiz_id = '" + quizList[i].ID + "';"
		rows, err := db.Query(querry)
		checkErr(err)
		for rows.Next() {
			question := Question{}
			err = rows.Scan(&question.ID, &question.Text)
			checkErr(err)
			questionList = append(questionList, question)
		}

		for j := 0; j < len(questionList); j++ {
			answerList := []Answer{}
			querry := "Select id, answer, correct from answer where question_id = '" + questionList[j].ID + "';"
			rows, err := db.Query(querry)
			checkErr(err)
			for rows.Next() {
				answer := Answer{}
				err = rows.Scan(&answer.ID, &answer.Text, &answer.Correct)
				checkErr(err)
				answerList = append(answerList, answer)
			}
			questionList[j].Answers = answerList
		}
		quizList[i].Questions = questionList
	}
	json.NewEncoder(w).Encode(quizList)
}

func Login(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")
	passwd := r.Header.Get("passwd")
	fmt.Print(email + "\n")
	fmt.Print(passwd + "\n")
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/quizdb")
	text := "Select exists (SELECT id FROM `users` where `email` = '" + email + "' and `password` = '" + passwd + "');"
	var exists bool
	rows, err := db.Query(text)
	checkErr(err)
	checkedUser := User{}
	for rows.Next() {
		err = rows.Scan(&exists)
		checkErr(err)
		if exists {
			fmt.Print("Email and password are correct \n")
			querry := "SELECT id, name FROM `users` where `email` = '" + email + "' and `password` = '" + passwd + "';"
			rows, err := db.Query(querry)

			for rows.Next() {
				err = rows.Scan(&checkedUser.ID, &checkedUser.Name)
				checkErr(err)
			}

		} else {
			fmt.Print("Email or password is wrong \n")
			checkedUser.ID = "0"
		}
	}
	json.NewEncoder(w).Encode(checkedUser)
}

func main() {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/quizdb")
	querry := "USE  quizdb;"
	_, err = db.Exec(querry)
	checkErr(err)

	router := mux.NewRouter()
	router.HandleFunc("/users", GetPerson).Methods("GET")
	router.HandleFunc("/checkuser", IfUserExists).Methods("GET")
	router.HandleFunc("/createuser", CreatePerson).Methods("POST")
	router.HandleFunc("/users/{id}", DeletePerson).Methods("DELETE")
	router.HandleFunc("/quiz", GetQuiz).Methods("GET")
	router.HandleFunc("/login", Login).Methods("GET")
	router.HandleFunc("/updatequiz", UpdateQuiz).Methods("POST")
	router.HandleFunc("/deletequiz", DeleteQuiz).Methods("POST")
	router.HandleFunc("/createquiz", CreateQuiz).Methods("GET")

	// Listening on port 3000
	http.ListenAndServe(":3000", router)
	fmt.Printf("We are running")
	//Starting func writer
}

func checkErr(err error) {
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
}
