package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/u0nel/accept"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		useragent := r.Header.Get("User-Agent")

		// we can't use a map[string]func(w, r), because maps aren't ordered
		types := []string{"text/plain", "text/html", "application/json", "image/png", "image/webp", "application/pdf"}

		requestedType := accept.ServeType(types, r.Header.Get("Accept"))
		w.Header().Add("Content-Type", requestedType)

		switch requestedType {
		case "text/plain":
			w.Write([]byte(r.Header.Get("User-Agent")))
		case "text/html":
			writeHtml(w, useragent, types)
		case "application/json":
			writeJson(w, useragent)
		case "image/png":
			writePng(w, useragent)
		case "image/webp":
			writeWebp(w, useragent)
		case "application/pdf":
			writePdf(w, useragent)
		default:
			http.Error(w, "Could not serve requested Type", http.StatusNotAcceptable)
		}
	})
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func writeHtml(w http.ResponseWriter, useragent string, types []string) {
	w.Write([]byte(`<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/water.css@2/out/water.css">
<h1>What's your user Agent?</h1>
It's
<pre><code>` + useragent + `</code></pre>
Unlike other websites like whatsmyuseragent.org, 
this website responds to the <code>Accept</code> Header.
It currently supports <code>` + strings.Join(types, "</code>, <code>") + `</code>
. If you want more, you can open an issue on <a href="https://github.com/u0nel/useragent">Github</a>
`))
}

func writeJson(w http.ResponseWriter, useragent string) {
	v := struct {
		UserAgent string `json:"user_agent"`
	}{useragent}
	json.NewEncoder(w).Encode(v)
}

func makeImage(useragent string) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 7*len(useragent), 30))
	x := 0
	y := 13
	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(useragent)
	return img
}

func writePng(w http.ResponseWriter, useragent string) {
	img := makeImage(useragent)
	png.Encode(w, img)
}

func writeWebp(w http.ResponseWriter, useragent string) {
	img := makeImage(useragent)
	options, _ := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	webp.Encode(w, img, options)
}

func writePdf(w http.ResponseWriter, useragent string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, useragent)
	pdf.Output(w)
}
