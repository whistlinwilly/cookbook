package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

const maxPage = 69

var recipeNameRegex = regexp.MustCompile(`<a class="post-card-permalink".*cooking\/(\D*)\/"`)
var imgRegex = regexp.MustCompile(`<img src="(.*)-.*(jpg|png)`)
var titleRegex = regexp.MustCompile(`<h1.*?>(.*)</h1>`)

func main() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	http.HandleFunc("/", recipeFactory(r))
	http.ListenAndServe(":8080", nil)
}

type QuickRecipe struct {
	Title, Image, URL string
}

func recipeFactory(r *rand.Rand) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		f, _ := ioutil.ReadFile("cookbook.tmpl")
		t, _ := template.New("cookbook").Parse(string(f))
		var title, image, url string
		for title, image, url = randomRecipe(r); title == "Post navigation"; title, image, url = randomRecipe(r) {
		}
		t.Execute(response, QuickRecipe{title, image, url})
		return
	}
}

func randomRecipe(r *rand.Rand) (title, image, url string) {
	n := r.Intn(maxPage)
	//fmt.Println(n)

	// Get random cooking page
	url = fmt.Sprintf("http://thepioneerwoman.com/cooking/page/%v/", n)
	body := makeRequestWithUserAgent(url)

	// Choose a random recipe from that page
	matches := recipeNameRegex.FindAllStringSubmatch(body, -1)
	m := r.Intn(len(matches))
	recipeName := matches[m][1]

	// Get the printable version of that recipe
	url = fmt.Sprintf("http://thepioneerwoman.com/cooking/%v/?printable_recipe", recipeName)
	body = makeRequestWithUserAgent(url)
	img := imgRegex.FindStringSubmatch(body)
	image = img[1] + "." + img[2]
	title = titleRegex.FindStringSubmatch(body)[1]
	return
}

func makeRequestWithUserAgent(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.110 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	resp.Body.Close()
	return string(body)
}
