package main

import (
	"fmt"
	"log"
	"lol/lol"
	"net/http"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var outputFileName = "text.txt"

func Download(w http.ResponseWriter, r *http.Request) {

	// Check if the request method is GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set the Content-Disposition header to indicate that it's an attachment.
	w.Header().Set("Content-Disposition", "attachment; filename="+outputFileName)

	// Set the Content-Type header based on the file type (e.g., text/plain).
	w.Header().Set("Content-Type", "text/plain")

	// Serve the file for download.
	http.ServeFile(w, r, outputFileName)
}

func saveText(saveString string) {
	outfile, err := os.Create(outputFileName)
	check(err)
	defer outfile.Close()

	// Replace "<br>" with newline characters "\n"
	saveString = strings.ReplaceAll(saveString, "<br>", "\n")

	_, err = outfile.WriteString(saveString)
	check(err)

	fmt.Printf("Output saved to %s\n", outputFileName)
}

// indexHandler handles HTTP requests to the root URL ("/").
func indexHandler(w http.ResponseWriter, r *http.Request) {
	bodyContent := "Add some text in the text box and enjoy :)"

	// Check if the request method is POST.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "400: Bad Request; Unable to parse form", http.StatusBadRequest)
			return
		}
		// Get the input text and style from the form.
		text := r.Form.Get("text")
		style := r.Form.Get("banner")
		// Generate ASCII art based on the input text and style.
		// Ensure the style file has a ".txt" extension.
		if !strings.Contains(style, ".txt") {
			style += ".txt"
		}
		fmt.Println(text)
		save := lol.AsciiArtTransform(text, style)

		store := "<br>"
		for i := 0; i < 8*len(save); i++ {
			store += save[i]
		}
		saveText(store)
		store = `<pre style="color: #333366; text-align: left;">` + store + "</pre>"
		bodyContent = store
		// bodyContent = save // Use the generated ASCII art as the response content.
	}

	// HTML template for the web page.
	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>ASCII Art Web</title>
		<link rel="stylesheet" href="/static/style.css">
	</head>
	<body>
		<div class="container">
			<form action="/" method="POST"> <!-- Change action to "/" to match the route -->
				<div>
					<h1>ASCII-art-web</h1>
					<h2>Style</h2>
				</div>
				<div class="white">
					<input type="radio" name="banner" value="standard">Standard
					<input type="radio" name="banner" value="shadow">Shadow
					<input type="radio" name="banner" value="thinkertoy">Thinkertoy
				</div>
				<div>
					<h3>Input text</h3>
				</div>
				<div>
					<textarea placeholder="Input text here" id="add-text" name="text" rows="4" cols="50"></textarea>
				</div>
				<div>
					<input type="submit" value="Generate">
				</div>
			</form>
			<div>
				<h3>Press download <strong>only</strong> after clicking generate!</h3>
			</div>
			<form action="/download" method="GET">
				<input type="submit" value="Download">
			</form>
		</div>
		<div class="box">
			<pre>%s</pre> <!-- Display the generated ASCII art -->
		</div>
		
	</body>
	</html>
	`
	// Send the HTML response with the generated ASCII art.
	fmt.Fprintf(w, tmpl, bodyContent)

}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// Register the indexHandler to handle requests at the root URL ("/").
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/download", Download)
	// Start an HTTP server on port 8080.
	fmt.Printf("Starting server at port 8080\n")
	http.ListenAndServe(":8080", nil)
}
