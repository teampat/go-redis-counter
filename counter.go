package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
	"gopkg.in/redis.v4"
)

func main() {
	http.HandleFunc("/count/", drawImage)
	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func drawImage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	redisIncr(id)
	log.Println(id)

	m := image.NewRGBA(image.Rect(0, 0, 1, 1))
	blue := color.RGBA{0, 0, 0, 0}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	var img image.Image = m
	writeImage(w, &img)
}

func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func redisIncr(key string) {

    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    err := client.Incr(key).Err()
    if err != nil {
        log.Println(err)
    }

}