//Test program code
package main

import (
	
	"net/http"
	"fmt"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"html/template"
	"sort"
	"flag"
	"log"

)

var (
	addr      = flag.String("addr", ":8000", "http service address")
	tmplname    = flag.String("tmplname", "post_services.html", "path to assets")
	homeTempl *template.Template
	
)

type JsonResponse struct {
    Rpc string `json:"jsonrpc"`
    Id int `json:"id"`
    Result []Response `json:"result"`
}

type Response struct {
	Country_code string `json:"country_code"`
	Slug string `json:"slug"`
	PostName string `json:"name"`
	
}

type Combined struct {
	Country_code   string
	Pairs []Response
}


type Responses []Response

func (slice Responses) Len() int {
    return len(slice)
}

func (slice Responses) Less(i, j int) bool {
    return slice[i].Country_code < slice[j].Country_code
}

func (slice Responses) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}


func prepareForWebOutput(o []Response) []Combined {
	
	type stringPair struct {
		
	str1 string
	str2 string
	
	}
	
	firstkey := o[0].Country_code
	ret := []Combined{}
	stringPairArray := []Response{}
	//stringPairElnumber := 0
	for elnum, v := range o {
		if v.Country_code != firstkey {
		
			ret = append(ret, Combined{o[elnum-1].Country_code, stringPairArray})
			
			firstkey = v.Country_code
			//fmt.Println("stringPairArray before clear", stringPairArray)
			stringPairArray = nil
			
			//fmt.Println("stringPairArray after clear", len(stringPairArray), stringPairArray)
			//fmt.Println("New first key", firstkey)
			//fmt.Println("Current v key", v.Country_code)
			
			newStringPair := Response{v.Country_code, v.Slug, v.PostName}
			stringPairArray = append(stringPairArray, newStringPair)
			
			if elnum == len(o)-1 {
			
				//fmt.Println("Reached last key: ", v.Country_code)
				
				ret = append(ret, Combined{v.Country_code, stringPairArray})
				return ret
				
			}
			
		
			
			
			
		} else {
			newStringPair :=  Response{v.Country_code, v.Slug, v.PostName}
			stringPairArray = append(stringPairArray, newStringPair)

			//fmt.Println("If String: ", stringPairArray)
			if elnum == len(o)-1 {
			
				//fmt.Println("Reached last key: ", v.Country_code)
				
				ret = append(ret, Combined{v.Country_code, stringPairArray})
				return ret
				
			}
			
			
			
			
			
		}
	}
	return ret
	
}



func apiFunc(w http.ResponseWriter, r *http.Request) {
		var apiKey string
   // var currentCountry string
	
	fmt.Println("Enter")
	apiKey = "e216ab38a1e2da5c0afdc4a2a8c0efc82b9743e8f9295b47277cdc071a60d1de5bfd9504e694e542"
	

	//resp,_ :=http.Get(http://gdeposylka.ru/api/v3/jsonrpc)
	
	
	
	
	url := "http://gdeposylka.ru/api/v3/jsonrpc"
    fmt.Println("URL:>", url)

    var jsonStr =[]byte(`
    
{
    "jsonrpc": "2.0",
    "method": "getCouriers",
    "params": { },
    "id": 1
}

`)

    
    
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Authorization-Token",apiKey)
    req.Header.Set("Content-Type", "application/json")
    
   

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, err := ioutil.ReadAll(resp.Body)
    
    if err != nil {
    	
    	log.Fatal("Cannot read JSON-RPC response body")
    }
    fmt.Println("response Body:", string(body))
    
    js := JsonResponse{}
    json.Unmarshal(body,&js)
    
 
    
    
    fmt.Println("Unsorted stuff: ",js.Result)
 
     t, err := template.ParseFiles(*tmplname)
     
     if err != nil {
     	
     	log.Fatal("Cannot read template file" + *tmplname)
     	
     } 
       
  	responses := Responses(js.Result)
    sort.Sort(responses)
	  fmt.Println("Sorted stuff: ",responses)
      fmt.Println("Sorted stuff: ",js.Result)
      
    
      fmt.Println ("Number of structs: ",len(js.Result))
      
      finalOutput := prepareForWebOutput(js.Result)
      
      fmt.Println("Final Output: ", finalOutput)
     
     //Циклом проходим по каждому элементу стракта. Если country_code в элементе i+1 == country_code в i, то 
     

     //js.Result
    
    
     
   w.Header().Add("Content Type", "text/html")
        t.Execute(w,finalOutput)
}

func main() {
	
	
	flag.Parse()
  
  
 
    http.HandleFunc("/", apiFunc)
    http.ListenAndServe(*addr,nil)
	
	
}

