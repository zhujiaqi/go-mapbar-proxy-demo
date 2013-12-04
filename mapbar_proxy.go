package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	//"mars"
	"math"
	"net/http"
)

func MapbarProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Header().Set("content-type", "text/plain;charset=utf-8")
		//read data in as whole and parse as json
		raw_data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "error", 500)
			return
		}
		parsed_data, err := simplejson.NewJson(raw_data)
		if err != nil {
			log.Println(err)
			http.Error(w, "not a valid json", 500)
			return
		}
		//set data and headers
		parsed_data.Set("version", "2.0.0")
		parsed_data.Set("host", "www.mapbar.com")
		if _, ok := parsed_data.CheckGet("radio_type"); !ok {
			parsed_data.Set("radio_type", "gsm")
		}
		if _, ok := parsed_data.CheckGet("location"); ok {
			parsed_data.Set("location", nil)
		}
		encoded_data, err := parsed_data.Encode()
		if err != nil {
			log.Println(err)
			panic("shit")
		}
		//make request to retrieve data from mapbar
		client := &http.Client{}
		req, _ := http.NewRequest("POST", "http://api.s.mapbar.com/position/getPosition.json", bytes.NewReader(encoded_data))
		req.Header.Set("Host", "api.s.mapbar.com")
		req.ContentLength = int64(len(encoded_data))
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			//TODO: better error
			log.Println(err)
			//http.Error(w, "error: request mapbar", 500)
			goto Ill_response
		}
		raw_resp_data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//TODO: better error
			log.Println(err)
			//http.Error(w, "error: read from response", 500)
			goto Ill_response
		}
		parsed_resp_data, err := simplejson.NewJson(raw_resp_data)
		if err != nil {
			//TODO: better error
			log.Println(err)
			//http.Error(w, "error: parse response json", 500)
			goto Ill_response
		}
		//result hack and slash
		result_map := make(map[string]interface{})
		for k, v := range parsed_resp_data.MustMap() {
			if k == "position" {
				result_map["location"] = v
			} else {
				result_map[k] = v
			}
		}
		if location, ok := result_map["location"]; ok {
			var latitude float64
			var longitude float64
			var accuracy int64
			latitude, _ = location.(map[string]interface{})["latitude"].(json.Number).Float64()
			longitude, _ = location.(map[string]interface{})["longitude"].(json.Number).Float64()
			accuracy, _ = location.(map[string]interface{})["accuracy"].(json.Number).Int64()
			if !parsed_data.Get("mars").MustBool() {
				if latitude == 0.0 && longitude == 0.0 && accuracy == 0 {
					log.Println("hack: mapbar returns [0, 0]")
					goto Ill_response
				}
				if math.Abs(latitude-39.904214) < 0.00005 && math.Abs(longitude-116.407413) < 0.00005 && accuracy > 15000 {
					log.Println("hack: remove this fake result")
					goto Ill_response
				}
				//latitude, longitude = mars.MarsToEarth(latitude, longitude)
				//log.Printf("new latlon: [%f, %f]\n", latitude, longitude)
				//location.(map[string]interface{})["latitude"] = latitude
				//location.(map[string]interface{})["longitude"] = longitude
			}
			if _, ok := location.(map[string]interface{})["address"]; ok {
				delete(location.(map[string]interface{}), "address")
			}
			log.Printf("response: [%f, %f] %d\n", latitude, longitude, accuracy)
		}
		//convert response data to json str
		encoded_resp_data, err := json.Marshal(result_map)
		if err != nil {
			log.Println(err)
			goto Ill_response
		}
		fmt.Fprint(w, string(encoded_resp_data))
	} else {
		http.Error(w, "Please POST your data", 405)
	}

	return

Ill_response:
	fmt.Fprint(w, "{}")
}

func main() {
	http.HandleFunc("/", MapbarProxy)
	err := http.ListenAndServe(":19999", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
