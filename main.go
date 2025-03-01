package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"maps"
	"net/http"
	"os"
	"slices"
)

type JupiterPriceResponse struct {
	Data map[string]Token `json:"data"`
}

type Token struct {
	Price string `json:"price"`
}

const JUPCA = "JUPyiwrYJFskUPiHa7hkeR8VUtAeFoSYbKedZNsDvCN"
const SOLCA = "So11111111111111111111111111111111111111112"
const BTCCA = "cbbtcf3aa214zXHbiAZQwf4122FBYbraNdFqgw4iMij"
const ETHCA = "EyrnrbE5ujd3HQG5PZd9MbECN9yaQrqc8pRwGtaLoyC"

var tokensMap = map[string]string{
	"JUP": JUPCA,
	"SOL": SOLCA,
	"BTC": BTCCA,
	"ETH": ETHCA,
}

func main() {
	if len(os.Args) > 1 {
		ca := os.Args[1]
		response, err := getJupiterPriceResponse(ca)
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}

		fmt.Println("Price from JUP: ", response.Price)
		return
	}

	prompt := promptui.Select{
		Label: "Choose token",
		Items: getMapKeys(),
	}

	_, result, _ := prompt.Run()

	response, err := getJupiterPriceResponse(tokensMap[result])

	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return
	}

	fmt.Println("Price from JUP: ", response.Price)
}

func getMapKeys() []string {
	return slices.Collect(maps.Keys(tokensMap))
}

func getJupiterPriceResponse(input string) (Token, error) {
	resp, err := http.Get("https://api.jup.ag/price/v2?ids=" + input)
	if err != nil {
		return Token{}, errors.New(err.Error())
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return Token{}, errors.New(err.Error())
	}

	var price JupiterPriceResponse
	err = json.Unmarshal(buf.Bytes(), &price)
	if err != nil {
		return Token{}, errors.New(err.Error())
	}

	if price.Data[input] == (Token{}) {
		return Token{}, errors.New("CA not found on jupiter")
	}

	return price.Data[input], nil
}
