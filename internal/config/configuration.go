package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	return write(*c)
}

func Read() (Config, error) {
	path, err := getConfigPath()
	if err != nil {
		fmt.Println("error al intentar obtener el home directory")
		return Config{}, err
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error al intentar abrir el archivo")
		return Config{}, err
	}
	defer file.Close()
	var cfg Config
	decoder := json.NewDecoder(file)             //Decoder acepta como parámetro un io.Reader, Unmarshal, un []byte
	if err := decoder.Decode(&cfg); err != nil { //todo lo que tenga un método Read, cumple la interfaz io.Reader
		fmt.Println("error al decodificar el archivo")
		return Config{}, err
	}
	return cfg, nil
}

func write(c Config) error {
	jsoned, err := json.Marshal(c) //Creas el json
	if err != nil {
		fmt.Printf("error al intentar crear el json")
		return err
	}
	path, err := getConfigPath() //Obtienes el path donde se debe actualizar el archivo con el usuario nuevo
	if err != nil {
		fmt.Printf("error al intentar escribir el json")
		return err
	}
	j, err := os.Create(path) //Creas el archivo/borras su contenido si ya existe
	if err != nil {
		fmt.Printf("error al intentar crear el archivo")
		return err
	}
	defer j.Close()                            //Cierras el archivo
	if _, err := j.Write(jsoned); err != nil { //Escribes el archivo. Write acepta []byte
		fmt.Printf("Error al escribir el json")
		return err
	}
	return nil
}

func getConfigPath() (string, error) {
	root, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("No se pudo obtener el home directory")
		return "", err
	}
	path := filepath.Join(root, filename) //Une los elementos en una sola ruta con le separador que utilice el sistema
	return path, nil
}
