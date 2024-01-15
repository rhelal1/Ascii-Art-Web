package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type PageData struct {
	Color string
	Text  []string
}

type PostRequest struct {
	input_text       string
	text_sytle       string
	user_color       string
	file_type        string
	download_request string
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/get-file", getFileHandler)
	http.ListenAndServe(":8080", nil)
}

func getFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment")
	http.ServeFile(w, r, "template/file")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// EARLY HANDLING AND EXIT
	if r.URL.Path == "/style.css" {
		http.ServeFile(w, r, "./template/style.css")
		return
	} else if r.URL.Path != "/" || (r.Method != "GET" && r.Method != "POST") {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "./template/404.html")
		return
	}
	if r.Method == "GET" {
		indexTemplate, _ := template.ParseFiles("template/index.html")
		err := indexTemplate.Execute(w, PageData{})
		if err != nil {
			fmt.Print(err)
		}
		return
	}
	rootHandlerPost(w, r)
}

func rootHandlerPost(w http.ResponseWriter, r *http.Request) {
	indexTemplate, _ := template.ParseFiles("template/index.html")

	if r.Header.Get("Content-type") != "application/x-www-form-urlencoded" || r.ContentLength > 2200 {
		return
	}

	// POST only from this point
	postRequest := getPostPramaters(r)
	is_request_vaild := postRequest != nil
	if !is_request_vaild {
		w.WriteHeader(http.StatusBadRequest)
		http.ServeFile(w, r, "./template/400.html")
		return
	}

	_, error := os.Stat(postRequest.text_sytle + ".txt")
	if os.IsNotExist(error) || len(postRequest.input_text) > 2000 {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "template/500.html")
		return
	}

	textInASCII := getAsciiString(postRequest.input_text, postRequest.text_sytle)
	pageData := PageData{
		Text:  textInASCII,
		Color: postRequest.user_color,
	}

	if postRequest.download_request == "yes" {
		WriteFile(postRequest.input_text, postRequest.text_sytle)
	}

	err := indexTemplate.Execute(w, pageData)
	if err != nil {
		fmt.Print(err)
	}
}

func getPostPramaters(r *http.Request) *PostRequest {

	input_text := r.FormValue("thetext")
	text_sytle := r.FormValue("chose")
	user_color := r.FormValue("color")
	file_type := r.FormValue("FileType")
	download_request := r.FormValue("download")
	vaild_input := CheckLetter(input_text)

	if !vaild_input || len(input_text) == 0 {
		return nil
	}

	return &PostRequest{
		input_text:       input_text,
		text_sytle:       text_sytle,
		user_color:       user_color,
		file_type:        file_type,
		download_request: download_request,
	}
}

func getAsciiString(text, filename string) []string {
	var Text []string
	WordsInArr := strings.Split(text, "\r\n")
	var Words []string
	for l := 0; l < len(WordsInArr); l++ {
		var Words [][]string
		Text1 := strings.ReplaceAll(WordsInArr[l], "\\t", "   ")
		if Text1 != "" {
			for j := 0; j < len(Text1); j++ {
				Words = append(Words, ReadLetter(Text1[j], filename))
			}
			for x := 0; x < 8; x++ {
				Lines := ""
				for n := 0; n < len(Words); n++ {
					Lines += Words[n][x]
				}
				Text = append(Text, Lines)
			}
		} else {
			Text = append(Text, "\n")
		}
	}
	return append(Words, strings.Join(Text, "\n"))
}

func ReadLetter(Text1 byte, fileName string) []string {
	var Letter []string
	ReadFile, _ := os.Open(fileName + ".txt")
	FileScanner := bufio.NewScanner(ReadFile)
	stop := 1
	i := 0
	letterLength := (int(Text1)-32)*9 + 2
	for FileScanner.Scan() {
		i++
		if i >= letterLength {
			stop++
			Letter = append(Letter, FileScanner.Text())
			if stop > 8 {
				break
			}
		}
	}
	ReadFile.Close()
	return Letter
}

func CheckLetter(s string) bool {
	WordsInArr := strings.Split(s, "\r\n")
	for l := 0; l < len(WordsInArr); l++ {
		for g := 0; g < len(WordsInArr[l]); g++ {
			if (WordsInArr[l][g] > 126 || WordsInArr[l][g] < 32) && WordsInArr[l][g] != 10 {
				return false
			}
		}
	}
	return true
}

func WriteFile(s, text_sytle string) {
	file, err := os.Create("template/file")
	if err != nil {
		fmt.Println("Error \n", err)
	} else {
		os.Chmod("template/file", 0600)
		WordsInArr := strings.Split(s, string(10))
		for l := 0; l < len(WordsInArr); l++ {
			var Words [][]string
			Text1 := strings.ReplaceAll(WordsInArr[l], "\\t", "   ")
			for j := 0; j < len(Text1); j++ {
				Words = append(Words, ReadLetter(Text1[j], text_sytle))
				if Text1[j] == 10 {
					Words = append(Words, []string{"\n"})
				}
			}
			if len(Words) != 0 {
				for w := 0; w < 8; w++ {
					for n := 0; n < len(Words); n++ {
						file.WriteString(Words[n][w])
					}
					if w+1 != 8 {
						file.WriteString("\n")
					}
				}
			}
			file.WriteString("\n")
		}
	}
	file.Close()
}
