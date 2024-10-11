package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	xtptest "github.com/dylibso/xtp-test-go"
	"github.com/extism/go-pdk"
)

//export test
func test() {
	topic := "input-topic"
	key := "key"

	xtptest.Group("series-1", func() {
		scanner := bufio.NewScanner(bytes.NewReader(xtptest.MockInputBytes()))
		for scanner.Scan() {
			date, price, volume, expected := parseLine(scanner.Text())

			input := Record{
				Topic: topic,
				Key:   key,
				Value: Order{
					Date:   date,
					Price:  price,
					Volume: volume,
				}}
			inputBytes, err := json.Marshal(input)
			if err != nil {
				pdk.Log(pdk.LogError, fmt.Sprintf("error %v", err))
				panic(err)
			}
			outputBytes := xtptest.CallBytes("transform", inputBytes)

			results := make([]Record, 1)
			json.Unmarshal(outputBytes, &results)
			given := results[0]

			xtptest.AssertEq(fmt.Sprintf("max(..., %.4f) = %.4f", price, given.Value.Price), expected, given.Value.Price)
		}
	})
}

func parseLine(txt string) (date time.Time, price float64, volume int64, expected float64) {
	values := strings.Split(txt, ",")
	date, _ = time.Parse(time.RFC3339, values[0])
	price, _ = strconv.ParseFloat(values[1], 64)
	volume, _ = strconv.ParseInt(values[2], 10, 64)
	expected, _ = strconv.ParseFloat(values[3], 64)
	return
}

// A key/value header pair.
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// An order from the market.
type Order struct {
	// Date/time of the order
	Date time.Time `json:"date"`
	// Closing price of the order
	Price float64 `json:"price"`
	// The volume of the order
	Volume int64 `json:"volume"`
}

// A plain key/value record.
type Record struct {
	Headers []Header `json:"headers"`
	Key     string   `json:"key"`
	Topic   string   `json:"topic"`
	// An order from the market.
	Value Order `json:"value"`
}

func main() {}
