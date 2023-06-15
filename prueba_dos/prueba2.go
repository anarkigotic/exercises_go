package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type User struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type Organization struct {
	Organization string `json:"organization"`
	Users        []User `json:"users"`
}

func main() {
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	result := []Organization{}

	// Iterar a partir de la segunda fila para omitir los títulos
	for i := 1; i < len(data); i++ {
		row := data[i]
		organization := row[0]
		username := row[1]
		role := row[2]

		// Verificar si la organización ya existe en el resultado
		orgIndex := -1
		for j, org := range result {
			if org.Organization == organization {
				orgIndex = j
				break
			}
		}

		// Si la organización no existe, se agrega al resultado
		if orgIndex == -1 {
			newOrg := Organization{
				Organization: organization,
				Users:        []User{},
			}
			result = append(result, newOrg)
			orgIndex = len(result) - 1
		}

		// Verificar si el usuario ya existe en la organización
		userIndex := -1
		for k, user := range result[orgIndex].Users {
			if user.Username == username {
				userIndex = k
				break
			}
		}

		// Si el usuario no existe, se agrega a la organización
		if userIndex == -1 {
			newUser := User{
				Username: username,
				Roles:    []string{},
			}
			result[orgIndex].Users = append(result[orgIndex].Users, newUser)
			userIndex = len(result[orgIndex].Users) - 1
		}

		// Agregar el rol al usuario
		result[orgIndex].Users[userIndex].Roles = append(result[orgIndex].Users[userIndex].Roles, role)
	}

	// Convertir result a JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}

	// Imprimir el resultado como JSON
	fmt.Println(string(jsonData))
}
