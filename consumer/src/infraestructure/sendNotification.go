package infraestructure

import (
	"bytes" 
	"context"
	"encoding/json"
	"fmt"
	"io" 
	"io/ioutil"
	"net/http"
	"os"
	"golang.org/x/oauth2/google"
)

type NotificationServiceImpl struct{}

func NewNotificationServiceImpl() *NotificationServiceImpl {
	return &NotificationServiceImpl{}
}

func (s *NotificationServiceImpl) SendNotification(token string) error {
	accessToken, err := getAccessToken()
	if err != nil {
		return fmt.Errorf("error obteniendo token de acceso: %v", err)
	}

	notif := map[string]interface{}{
		"message": map[string]interface{}{
			"token": token,
			"notification": map[string]string{
				"title": "Registro exitoso",
				"body":  "¡Tu token fue guardado exitosamente!",
			},
			"data": map[string]string{
				"status": "success",
			},
		},
	}

	jsonData, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("error generando JSON de la notificación: %v", err)
	}

	IDcliente := os.Getenv("ID_CLIENTE")
	if IDcliente == "" {
		return fmt.Errorf("la variable de entorno 'ID_CLIENTE' no está configurada")
	}

	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", IDcliente)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando solicitud HTTP: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando notificación: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Respuesta FCM:", string(body))
	return nil
}

func (s *NotificationServiceImpl) SendFCMNotification(token, header, description, image string) error {
	accessToken, err := getAccessToken()
	if err != nil {
		return fmt.Errorf("error obteniendo token de acceso: %v", err)
	}

	notif := map[string]interface{}{
		"message": map[string]interface{}{
			"token": token,
			"notification": map[string]string{
				"title": header,
				"body":  description,
			},
			"data": map[string]string{
				"image": image,
			},
		},
	}

	jsonData, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("error generando JSON de la notificación: %v", err)
	}

	IDcliente := os.Getenv("ID_CLIENTE")
	if IDcliente == "" {
		return fmt.Errorf("la variable de entorno 'ID_CLIENTE' no está configurada")
	}

	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", IDcliente)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando solicitud HTTP: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando notificación: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Respuesta de FCM:", string(body))
	return nil
}

func getAccessToken() (string, error) {
	data, err := ioutil.ReadFile("./multi.json")
	if err != nil {
		return "", fmt.Errorf("error leyendo archivo de cuenta de servicio: %v", err)
	}

	config, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return "", fmt.Errorf("error creando configuración JWT: %v", err)
	}

	token, err := config.TokenSource(context.Background()).Token()
	if err != nil {
		return "", fmt.Errorf("error obteniendo token de acceso: %v", err)
	}

	return token.AccessToken, nil
}