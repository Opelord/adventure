package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var fileName = flag.String("f", "gopher.json", "Specify the name of the story file")

// tpl is a template of how the story page will look like.
// The body of the page will show story and beneath it will be hyperlinked options with links to the next chapter
// In case of no further options "THE END" will appear.
const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<div>{{ range .Story }}{{ . }}</br>{{end}}</div>
		<form>
		{{ range .Options }}<a href="/{{ .Arc }}">{{ .Text }}</a><br>{{else}}<div><strong>THE END</strong></div>{{end}}
		</form>
	</body>
</html>
`

// Options and place are structures that are used for parsing structured .json file
type Options struct{
	Text string `json:"text"`
	Arc string `json:"arc"`
}

type place struct{
	Title string `json:"title"`
	Story []string `json:"story"`
	Options []Options `json:"options"`
}

// main has only two functions responsible for handling and listening on a given port
func main() {
	http.HandleFunc("/", tempHandler)
	http.ListenAndServe(":8080", nil)
}

// tempHandler is serving properly parsed template with provided data on the server
func tempHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	t, err := template.New("webpage").Parse(tpl)
	if err != nil{
		fmt.Println("Couldn't parse template")
		fmt.Printf("%e\n", err)
	}
	d := getData()
	t.Execute(w, d[title])
}

// getData is parsing .json file and returning it in a form of a map[string]place
func getData() map[string]place	{
	flag.Parse()
	f, err := os.Open(*fileName)
	if err != nil {
		fmt.Println("Couldn't open file")
	}
	defer f.Close()

	jsonData, err := ioutil.ReadAll(f)
	var data map[string]place
	err = json.Unmarshal(jsonData, &data)
	if err != nil{
		fmt.Printf("%e\n", err)
	}
	return data
}
