package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	authusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/auth/usecase"
)

// Register є HTTP-обробником (handler), який приймає запити на реєстрацію.
// Він пов'язує зовнішній HTTP-інтерфейс із внутрішньою бізнес-логікою (usecase).
func Register(svc authusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input authusecase.RegisterInput

		// Декодування вхідного JSON-пакета з тіла HTTP-запиту у Go-структуру.
		// Якщо JSON некоректний (наприклад, пропущена дужка), повертаємо помилку 400.
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Ініціалізація контексту із обмеженням часу виконання у 30 секунд.
		// Це захищає сервер від «зависання» запиту та витоку ресурсів (пам'яті, з'єднань з БД).
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel() // Гарантоване звільнення ресурсів контексту після завершення функції

		// Передача управління на шар бізнес-логіки (usecase) для обробки реєстрації.
		resp, err := svc.Register(ctx, input)
		if err != nil {
			// 4. Диференціація (розділення) помилок для коректної відповіді клієнту.
			// Якщо пошта вже зайнята або введені дані не пройшли валідацію — це помилка клієнта (400 Bad Request).
			if errors.Is(err, authusecase.ErrEmailExists) || errors.Is(err, authusecase.ErrInvalidInput) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Перехоплення всіх інших непередбачуваних проблем (наприклад, збій бази даних).
			// У такому випадку клієнту повертається загальний статус 500 Internal Server Error.
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Формування успішної HTTP-відповіді.
		// Встановлюємо заголовок, що дані віддаються у форматі JSON, та статус 201 Created.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Серіалізація (кодування) об'єкта відповіді з Go-структури назад у JSON-формат для клієнта.
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func Login(svc authusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input authusecase.LoginInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		resp, err := svc.Login(ctx, input)
		if err != nil {
			if errors.Is(err, authusecase.ErrInvalidCredentials) || errors.Is(err, authusecase.ErrInvalidInput) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func Validate(svc authusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			http.Error(w, "Missing or invalid Authorization header", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		resp, err := svc.ValidateToken(ctx, token)
		if err != nil {
			if errors.Is(err, authusecase.ErrInvalidToken) || errors.Is(err, authusecase.ErrTokenExpired) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
