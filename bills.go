package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
)

func main() {
	//Inicia opciones de servicios
	var serviceOptions ServiceOptions
	//prompt de escogido de servicios
	//qs son las opciones a escoger, sOptions el struct que guarda lo escogido
	err := ServicesPrompt(&serviceOptions)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Escogiste %s. \n", serviceOptions.Chosen)
	if serviceOptions.Chosen == "luz" {
		LuzService(serviceOptions)
	}
	//token := getLumaToken()
	//getBills(token)
}
func LuzService(options ServiceOptions) {
	fmt.Println("Escogiste luz")
	err := survey.Ask(luzqs, &options)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Escogiste %s. \n", options.Chosen)

	if options.Chosen == "Historial" {
		var token string
		getLumaToken(&token)
		fmt.Printf("Token %s. \n", token)
		getBills(token)
	}

}

func getLumaToken(token *string) {
	httposturl := "https://api.miluma.lumapr.com/miluma-api/auth"
	fmt.Println("HTTP POST URL", httposturl)
	var jsonData = []byte(`{"password": "PLACEHOLDER","username": "PLACEHOLDER"}`)
	request, error := http.NewRequest("POST", httposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	response, error := client.Do(request)

	if error != nil {
		panic(error)
	}
	body, _ := ioutil.ReadAll(response.Body)
	var res PostResponse
	json.Unmarshal([]byte(string(body)), &res)
	fmt.Printf("Token: %s", res.Data.Token)
	fmt.Printf("Message: %s", res.Message)
	*token = res.Data.Token
}

func getBills(token string) {
	url := "https://api.miluma.lumapr.com/miluma-bill-api/api/bill/history?accId={AccountID}"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var bills Bills
	json.Unmarshal([]byte(string(body)), &bills)
}
func ServicesPrompt(serviceOptions *ServiceOptions) error {
	return survey.Ask(qs, serviceOptions)
}

type PostResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Message string `json:"message"`
}

type Bills struct {
	AccID string `json:"accId"`
	Bills []struct {
		BillID     string `json:"billId"`
		BillDate   string `json:"billDate"`
		BillAmount string `json:"billAmount"`
	} `json:"bills"`
}

type ServiceOptions struct {
	Chosen string `survey:"option"` //'option' es el field name para survey questions
}

var qs = []*survey.Question{
	{
		Name: "option", //field name que buscara en el struct y guarda el value
		Prompt: &survey.Select{
			Message: "Escoge:",
			Options: []string{"luz", "agua", "internet/telefonica"},
			Default: "luz",
		},
	},
}

var luzqs = []*survey.Question{
	{
		Name: "option", //field name que buscara en el struct y guarda el value
		Prompt: &survey.Select{
			Message: "Escoge:",
			Options: []string{"Historial", "Pagar"},
			Default: "Historial",
		},
	},
}
