package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)
//Foods struct which contains an array of foods
type Foods struct {
	Foods []Food `json:"foods"`
}
//Food struct which contains details about the dish
type Food struct {
	Id               int              `json:"id"`
	Name             string           `json:"name"`
	PreparationTime  int              `json:"preparation-time"`
	Complexity       int              `json:"complexity"`
	CookingApparatus string           `json:"cooking-apparatus"`
}
//Order struct which contains details about the generated order
type Order struct {
	Id         int    `json:"id"`
	Items      []int  `json:"items"`
	Priority   int    `json:"priority"`
	MaxWait    int    `json:"max-wait"`
	PickUpTime int    `json:"pick-up-time"`
}

func genRandNum(min, max int64) int64 {
	// calculate the max we will be using
	bg := big.NewInt(max - min)

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}

	// add n to min to support the passed in range
	return n.Int64() + min
}

func generateItems() []int{
	n := int(genRandNum(1, 10))
	var items = make([]int, n)
	for i := 0; i < n; i++{
		items[i] = int(genRandNum(1, 10))
	}
	return items
}

func getMaxWait(foods Foods) int{
	maxTime := 0
	for i := 0; i < len(foods.Foods); i++{
		if foods.Foods[i].PreparationTime > maxTime{
			maxTime = foods.Foods[i].PreparationTime
		}
	}
	maxWait := float32(maxTime) * 1.3
	return int(maxWait)
}
func getUnixTimestamp() int {
	now := time.Now()
	sec := now.Unix()
	return int(sec)
}

func createOrder(foods Foods) []byte{

	order := &Order{Id : int(genRandNum(1, 100)),
		Items : generateItems(),
		Priority: int(genRandNum(1, 5)),
		MaxWait: getMaxWait(foods),
		PickUpTime: getUnixTimestamp(),
	}
	ord, err := json.Marshal(order)
	if err != nil{
		fmt.Printf("Error: %s", err)
	}
	return ord
}

func makeRequest (ord []byte){
	req, err := http.NewRequest("POST", "http://localhost:8080/dininghall", bytes.NewBuffer(ord))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Skip the error")
	} else {
		defer resp.Body.Close()
		fmt.Println(string(ord))
		fmt.Println("Request sent")
	}
}

func servePage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var order Order
	err := decoder.Decode(&order)
	if err != nil {
		panic(err)
	}
	fmt.Println("Request Handled")
	log.Println(order)
}

func waiter(foods Foods){
	order := createOrder(foods)
	makeRequest(order)
}
func main() {
	jsonFile, err := os.Open("../config/foods.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	//read our opened jsonFile as a byte array
	byteValue, _ := ioutil.ReadAll(jsonFile)

	//we initialize our Foods array
	var foods Foods

	//we unmarshal our byteArray which contains our
	//jsonFile's content info 'foods' which we defined above
	json.Unmarshal(byteValue, &foods)


	go func() {
		for {
			go func(){
				waiter(foods)
			}()
			time.Sleep(time.Second)
		}
	}()

	http.HandleFunc("/dininghall", servePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}