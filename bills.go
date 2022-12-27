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

	if serviceOptions.Chosen == "luz" {
		LuzService(serviceOptions)
	}

	if serviceOptions.Chosen == "internet/telefonica" {
		LibertyService(serviceOptions)
	}
}
func LuzService(options ServiceOptions) {
	fmt.Println("Escogiste luz")
	err := survey.Ask(luzqs, &options)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Escogiste %s. \n", options.Chosen)
	var token string
	if options.Chosen == "Historial" {
		getLumaToken(&token)
		fmt.Printf("Token %s. \n", token)
		getLumaBills(token)
	}

}
func LibertyService(options ServiceOptions) {
	httposturl := "https://mi.libertypr.com/Auth/LoginCDCAPI"
	fmt.Println("HTTP POST URL", httposturl)
	var jsonData = []byte(`{"Password": "PLACEHOLDER","Username": "PLACEHOLDER"}`)
	request, error := http.NewRequest("POST", httposturl, bytes.NewBuffer(jsonData))
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	body, _ := ioutil.ReadAll(response.Body)
	var res LibertyResponse
	json.Unmarshal([]byte(string(body)), &res)
	if !res.Success {
		for i := 0; i < len(res.Errors); i++ {
			fmt.Println("Error: ", res.Errors[i])
		}
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
	var res LumaResponse
	json.Unmarshal([]byte(string(body)), &res)
	fmt.Printf("Token: %s", res.Data.Token)
	fmt.Printf("Message: %s", res.Message)
	*token = res.Data.Token
}

func getLumaBills(token string) {
	url := "https://api.miluma.lumapr.com/miluma-bill-api/api/bill/history?accId={AccountID}"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{}
	response, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	body, _ := ioutil.ReadAll(response.Body)
	var bills Bills
	json.Unmarshal([]byte(string(body)), &bills)
}
func ServicesPrompt(serviceOptions *ServiceOptions) error {
	return survey.Ask(qs, serviceOptions)
}

type LumaResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Message string `json:"message"`
}
type LibertyResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
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
			Message: "Escoge servicio:",
			Options: []string{"luz", "internet/telefonica"},
			Default: "luz",
		},
	},
}

var luzqs = []*survey.Question{
	{
		Name: "option", //field name que buscara en el struct y guarda el value
		Prompt: &survey.Select{
			Message: "Escoge accion:",
			Options: []string{"Historial", "Pagar"},
			Default: "Historial",
		},
	},
}
