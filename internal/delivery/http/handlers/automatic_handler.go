package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AutomaticHandler struct {

}

func NewAutomaticHandler() *AutomaticHandler {
	return &AutomaticHandler{}
}

func (h *AutomaticHandler) Sms(c *gin.Context) {
    // Используем map[string]interface{} для любого JSON
    var body map[string]interface{}

    // Парсим JSON из тела запроса
    if err := c.BindJSON(&body); err != nil {
        log.Println("Ошибка при парсинге JSON:", err)
        c.JSON(400, gin.H{"error": "Invalid JSON"})
        return
    }

    // Логируем полученный реквест
    log.Println("Получен реквест:", body)

    // Можно дополнительно вывести тип каждого поля
    for k, v := range body {
        log.Printf("Поле: %s, Тип: %T, Значение: %v\n", k, v, v)
    }

    // Отправляем ответ
    c.JSON(200, gin.H{
        "status": "ok",
        "received": body,
    })
}


func (h *AutomaticHandler) Live(c *gin.Context) {
	// Используем map[string]interface{} для любого JSON
	var body map[string]interface{}
	// Парсим JSON из тела запроса
	if err := c.BindJSON(&body); err != nil {
		log.Println("Ошибка при парсинге JSON:", err)
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	// Логируем полученный реквест
	log.Println("Получен реквест:", body)

	// Можно дополнительно вывести тип каждого поля
	for k, v := range body {
		log.Printf("Поле: %s, Тип: %T, Значение: %v\n", k, v, v)
	}

	c.JSON(http.StatusOK, gin.H{
		"alive": true,
		"received": body,
	})
}

func (h *AutomaticHandler) Auth(c *gin.Context) {
	// Используем map[string]interface{} для любого JSON
	var body map[string]interface{}
	// Парсим JSON из тела запроса
	if err := c.BindJSON(&body); err != nil {
		log.Println("Ошибка при парсинге JSON:", err)
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	// Логируем полученный реквест
	log.Println("Получен реквест:", body)

	// Можно дополнительно вывести тип каждого поля
	for k, v := range body {
		log.Printf("Поле: %s, Тип: %T, Значение: %v\n", k, v, v)
	}

	// Отправляем ответ
	c.JSON(200, gin.H{
		"status": "ok",
		"received": body,
	})
}