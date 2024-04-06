package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type myElement struct {
	Name    string
	Surname string
	Id      string
}

var DATA = make(map[string]myElement)
var DATAFILE = "/tmp/dataFile.gob"

func save() error {
	fmt.Println("Saving", DATAFILE)
	err := os.Remove(DATAFILE)
	if err != nil {
		fmt.Println(err)
	}
	saveTo, err := os.Create(DATAFILE)
	if err != nil {
		fmt.Println("Cannot create", DATAFILE)
		return err
	}
	defer saveTo.Close()
	encoder := gob.NewEncoder(saveTo)
	err = encoder.Encode(DATA)
	if err != nil {
		fmt.Println("Cannot save to", DATAFILE)
		return err
	}
	return nil
}

func load() error {
	fmt.Println("Loading", DATAFILE)
	loadFrom, err := os.Open(DATAFILE)
	defer loadFrom.Close()
	if err != nil {
		fmt.Println("Empty")
		return err
	}

	decoder := gob.NewDecoder(loadFrom)
	decoder.Decode(&DATA)
	return nil
}

func Insert(k string, n myElement) bool {
	if k == "" {
		return false
	}
	if LOOKUP(k) == nil {
		DATA[k] = n
		return true
	}
	return false
}

func Delete(k string) bool {
	if LOOKUP(k) != nil {
		delete(DATA, k)
		return true
	}
	return false
}

func LOOKUP(k string) *myElement {
	_, ok := DATA[k]
	if ok {
		n := DATA[k]
		return &n
	} else {
		return nil
	}
}

func Update(k string, n myElement) bool {
	_, ok := DATA[k]
	if !ok {
		return false
	} else {
		DATA[k] = n
	}
	return true
}

func PRINT() {
	for k, v := range DATA {
		fmt.Println("key: %s value: %v\n", k, v)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving", r.Host, "for", r.URL.Path)
	myT := template.Must(template.ParseGlob("home.gohtml"))
	myT.ExecuteTemplate(w, "home.gohtml", nil)
}

func listAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listing the contents of the KV store!")
	fmt.Fprintf(w, "<a href=\"/\" style=\"margin-right: 20px;\"Home sweet home!</a>")
	fmt.Fprintf(w, "<a href=\"/list\" style=\"margin-right: 20px;\"List all elements</a>")
	fmt.Fprintf(w, "<a href=\"/change\" style=\"margin-right: 20px;\"Change an element</a>")
	fmt.Fprintf(w, "<a href=\"/insert\" style=\"margin-right: 20px;\"Insert an element</a>")
	fmt.Fprintf(w, "<a href=\"/delete\" style=\"margin-right: 20px;\"Delete an element</a>")

	fmt.Fprintf(w, "<h1>The contents of the KV store are:</h1>")
	fmt.Fprintf(w, "<ul>")
	for k, v := range DATA {
		fmt.Fprintf(w, "<li>")
		fmt.Fprintf(w, "<strong>%s</strong> with value: %v\n", k, v)
		fmt.Fprintf(w, "</li>")
	}
	fmt.Fprintf(w, "</ul>")
}

func changeElement(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Change an element of the KV store!")
	tmpl := template.Must(template.ParseFiles("update.gohtml"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	key := r.FormValue("key")
	n := myElement{
		Name:    r.FormValue("name"),
		Surname: r.FormValue("surname"),
		Id:      r.FormValue("id"),
	}
	if !Update(key, n) {
		fmt.Println("Update operation failed!")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	} else {
		err := save()
		if err != nil {
			fmt.Println(err)
			return
		}
		tmpl.Execute(w, struct{ Struct bool }{true})
	}
}

func insertElement(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inserting an element of the KV store!")
	tmpl := template.Must(template.ParseFiles("insert.gohtml"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	key := r.FormValue("key")
	n := myElement{
		Name:    r.FormValue("name"),
		Surname: r.FormValue("surname"),
		Id:      r.FormValue("id"),
	}
	if !Insert(key, n) {
		fmt.Println("Add operation failed!")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	} else {
		err := save()
		if err != nil {
			fmt.Println(err)
			return
		}
		tmpl.Execute(w, struct{ Struct bool }{true})
	}
}

func deleteElement(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Deleting an element of the KV store!")
	tmpl := template.Must(template.ParseFiles("delete.gohtml"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	key := r.FormValue("key")
	if !Delete(key) {
		fmt.Println("Delete operation failed!")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	} else {
		err := save()
		if err != nil {
			fmt.Println(err)
			return
		}
		tmpl.Execute(w, struct{ Struct bool }{true})
	}
}

func errorPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving", r.Host, "for", r.URL.Path)
	fmt.Fprintf(w, "<h1>Key error</h1>")
	fmt.Fprintf(w, "<a href=\"/\" style=\"margin-right: 20px;\"Home sweet home!</a>")
}

func main() {
	err := load()
	if err != nil {
		fmt.Println(err)
	}
	PORT := ":8081"
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Using default port number: ", PORT)
	} else {
		PORT = ":" + arguments[1]
	}

	http.HandleFunc("/", homePage)
	http.HandleFunc("/list", listAll)
	http.HandleFunc("/change", changeElement)
	http.HandleFunc("/insert", insertElement)
	http.HandleFunc("/delete", deleteElement)
	http.HandleFunc("/error", errorPage)
	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
	}

}
