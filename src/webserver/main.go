package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type CanvasResult struct {
	TotalPages  string `xml:"TotalPages"`
	CurrentPage string `xml:"CurrentPage"`
	Submissions []struct {
		Id string `xml:"Id,attr"`
	} `xml:"Submissions>Submission"`
}

func displayHtml(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprint(w, htmlBody()) // send data to client side
	getSubmissions(w)
	return
}

func getSubmissions(w http.ResponseWriter) {
	xmlBytes, err := gocanvasSubmissionsCall()
	if err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		var result CanvasResult
		xml.Unmarshal(xmlBytes, &result)
		fmt.Fprint(w, "<ul>")
		for i := 0; i < len(result.Submissions); i++ {
			fmt.Fprint(w, "<li>")
			fmt.Fprint(w, result.Submissions[i].Id)
			fmt.Fprint(w, "</li>")
		}

	}
}

func htmlBody() string {
	html := "<h1>Submissions List</h1>"
	return html
}

func gocanvasSubmissionsCall() ([]byte, error) {
	response, err := http.Get("https://demo.gocanvas.com/apiv2/submissions.xml?username=[username]&password=[password]&form_id=603606")
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}
	fmt.Println(data)
	return data, nil
}

func main() {
	http.HandleFunc("/", displayHtml)        // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
