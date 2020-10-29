package goseaweed

import (
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestSeaweed(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	fs := NewSeaweedFs("http://192.168.88.11:8888", time.Second*10)
	file, err := ioutil.ReadFile("./README/82274b7ec75f6f1fd86641f64dc82958.png")
	if err != nil {
		log.Fatalln(err)
	}
	if err := fs.PutObject("82274b7ec75f6f1fd86641f64dc82958.png", file); err != nil {
		log.Fatalln(err)
	}
	object, err := fs.GetObject("82274b7ec75f6f1fd86641f64dc82958.png")
	if err != nil {
		log.Fatalln(err)
	}
	ioutil.WriteFile("xxx.png",object,000666)
}
