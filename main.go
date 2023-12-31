package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
)

type WeatherData struct {
	Date      time.Time        `json:"date"`
	Forecasts []WeatherForecast `json:"forecasts"`
}

type WeatherForecast struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}

var client *mongo.Client
var apiKey = "sua_key_api"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/previsao", GetWeatherForecast).Methods("GET")

	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	addWeatherDataToDatabase()

	http.Handle("/", r)
	log.Println("Servidor iniciado na porta 8080")
	http.ListenAndServe(":8080", r)
}


func fetchWeatherData(city string) (WeatherForecast, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherForecast{}, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return WeatherForecast{}, err
	}

	if data["cod"].(float64) != 200 {
		errorMessage := fmt.Sprintf("Erro ao buscar dados de previsão do tempo: %s", data["message"].(string))
		return WeatherForecast{}, fmt.Errorf(errorMessage)
	}

	temperature := data["main"].(map[string]interface{})["temp"].(float64)
	description := data["weather"].([]interface{})[0].(map[string]interface{})["description"].(string)

	return WeatherForecast{
		Location:    city,
		Temperature: temperature,
		Description: description,
	}, nil
}

func addWeatherDataToDatabase() {
	cities := []string{"Sao Paulo", "Rio de Janeiro", "Brasília", "Florianópolis", "Campo Grande", "Fortaleza"}
	forecasts := []WeatherForecast{}

	for _, city := range cities {
		data, err := fetchWeatherData(city)
		if err != nil {
			log.Printf("Erro ao buscar dados de previsão do tempo para %s: %v\n", city, err)
		} else {
			forecasts = append(forecasts, data)
		}
	}

	weatherData := WeatherData{
		Date:      time.Now(),
		Forecasts: forecasts,
	}

	collection := client.Database("previsao-do-tempo").Collection("previsoes")

	opts := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(context.TODO(), bson.M{}, weatherData, opts)
	if err != nil {
		log.Printf("Erro ao atualizar/inserir dados de previsão do tempo: %v\n", err)
	} else {
		log.Println("Dados de previsão do tempo atualizados/inseridos com sucesso")
	}
}


func GetWeatherForecast(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("previsao-do-tempo").Collection("previsoes")

	var result WeatherData
	err := collection.FindOne(context.TODO(), bson.M{}).Decode(&result)
	if err != nil {
		log.Println("Erro ao buscar previsão do tempo:", err)
		http.Error(w, "Erro ao buscar previsão do tempo", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Println("Erro ao carregar o template:", err)
		http.Error(w, "Erro ao carregar a página", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, result)
	if err != nil {
		log.Println("Erro ao renderizar o template:", err)
		http.Error(w, "Erro ao renderizar a página", http.StatusInternalServerError)
	}
}
