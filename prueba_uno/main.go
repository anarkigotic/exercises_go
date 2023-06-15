package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response []Item

type Item struct {
	ClientID int     `json:"clientId"`
	Nombre   string  `json:"nombre"`
	Compro   bool    `json:"compro"`
	TDC      string  `json:"tdc,omitempty"`
	Monto    float64 `json:"monto,omitempty"`
	Date     string  `json:"date"`
}

type Result struct {
	Total         float64            `json:"total"`
	ComprasPorTDC map[string]float64 `json:"comprasPorTDC"`
	NoCompraron   int                `json:"nocompraron"`
	CompraMasAlta float64            `json:"compraMasAlta"`
}

type APIClient struct {
	client http.Client
}

func main() {
	apiClient := &APIClient{
		client: http.Client{},
	}

	http.HandleFunc("/resumen/", handleResumen(apiClient))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}

}

func fetchDataForDays(fechaInicial string, numDias int) (Response, error) {
	layout := "2006-01-02" // Formato de fecha
	startDate, err := time.Parse(layout, fechaInicial)
	if err != nil {
		return nil, fmt.Errorf("error al analizar la fecha inicial: %v", err)
	}

	fechas := make(Response, 0)

	resultados := make(chan []byte)

	for i := 0; i < numDias; i++ {
		go fetchDataAsync(startDate, i, layout, resultados)
	}

	for i := 0; i < numDias; i++ {
		body := <-resultados
		if body != nil {
			var response Response
			err = json.Unmarshal(body, &response)
			if err != nil {
				fmt.Println("Error al analizar la respuesta:", err)
				continue
			}
			fechas = append(fechas, response...)
		}
	}

	close(resultados)

	return fechas, nil
}

func fetchDataAsync(startDate time.Time, i int, layout string, resultados chan []byte) {
	fecha := startDate.AddDate(0, 0, i).Format(layout)
	urlBase := "https://apirecruit-gjvkhl2c6a-uc.a.run.app/compras/"
	url := urlBase + fecha
	body, err := fetchData(url)
	if err != nil {
		fmt.Println("Error al obtener datos:", err)
		resultados <- nil
		return
	}
	resultados <- body
}

func fetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func processResponse(items []Item) Result {
	result := Result{
		Total:         0,
		ComprasPorTDC: make(map[string]float64),
		NoCompraron:   0,
		CompraMasAlta: 0,
	}

	for _, item := range items {
		result.Total += item.Monto
		if item.Compro {
			result.ComprasPorTDC[item.TDC] += item.Monto
		} else {
			result.NoCompraron++
		}
		if item.Monto > result.CompraMasAlta {
			result.CompraMasAlta = item.Monto
		}
	}

	return result
}

func handleResumen(apiClient *APIClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := strings.Split(r.URL.Path, "/")
		if len(params) != 3 {
			http.Error(w, "Ruta inválida", http.StatusBadRequest)
			return
		}

		date := params[2]
		dias := r.URL.Query().Get("dias")

		daysToAdd := 0

		if dias != "" {
			var err error
			daysToAdd, err = strconv.Atoi(dias)
			if err != nil {
				http.Error(w, "Error en el parámetro 'dias'", http.StatusBadRequest)
				return
			}
		}

		fechas, err := fetchDataForDays(date, daysToAdd)
		if err != nil {
			fmt.Println("Error al obtener las fechas:", err)
			return
		}
		data := processResponse(fechas)
		resultJSON, err := json.Marshal(data)

		fmt.Println(string(resultJSON))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	}
}
