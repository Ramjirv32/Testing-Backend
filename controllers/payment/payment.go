package payment

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v3"

	"backend/utils"
)

// orderCounter drives round-robin gateway selection.
// Odd  => Razorpay
// Even => Cashfree
var orderCounter uint64

// testAmountINR is charged while testing. Change to real amount when going live.
const testAmountINR = 1.0

// ---- Request / Response types -----------------------------------------------

type CreateOrderReq struct {
	Amount        float64 `json:"amount"`
	BookingRef    string  `json:"booking_ref"`
	CustomerName  string  `json:"customer_name"`
	CustomerEmail string  `json:"customer_email"`
	CustomerPhone string  `json:"customer_phone"`
}

type VerifyReq struct {
	Gateway   string `json:"gateway"`
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Signature string `json:"signature"`
}

// ---- CreateOrder ------------------------------------------------------------

func CreateOrder(c fiber.Ctx) error {
	var req CreateOrderReq
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	// Override with test amount for now
	amount := testAmountINR

	count := atomic.AddUint64(&orderCounter, 1)
	if count%2 == 1 {
		return createRazorpayOrder(c, amount, req)
	}
	return createCashfreeOrder(c, amount, req)
}

// ---- Razorpay ---------------------------------------------------------------

func createRazorpayOrder(c fiber.Ctx, amountINR float64, req CreateOrderReq) error {
	keyID := os.Getenv("RAZORPAY_KEY_ID")
	keySecret := os.Getenv("RAZORPAY_KEY_SECRET")

	amountPaise := int64(amountINR * 100)
	receipt := fmt.Sprintf("tp_%d", time.Now().UnixMilli())

	payload, _ := json.Marshal(map[string]interface{}{
		"amount":   amountPaise,
		"currency": "INR",
		"receipt":  receipt,
	})

	httpReq, err := http.NewRequest("POST", "https://api.razorpay.com/v1/orders", bytes.NewReader(payload))
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to prepare Razorpay request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.SetBasicAuth(keyID, keySecret)

	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(httpReq)
	if err != nil {
		return utils.ErrorResponse(c, 502, "Razorpay unreachable")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode != 200 {
		msg := "Razorpay order creation failed"
		if errObj, ok := result["error"].(map[string]interface{}); ok {
			if d, ok := errObj["description"].(string); ok {
				msg = d
			}
		}
		return utils.ErrorResponse(c, 502, msg)
	}

	orderID, _ := result["id"].(string)
	return c.JSON(fiber.Map{
		"status": 200,
		"data": fiber.Map{
			"gateway":  "razorpay",
			"order_id": orderID,
			"key_id":   keyID,
			"amount":   amountINR,
			"currency": "INR",
		},
	})
}

// ---- Cashfree ---------------------------------------------------------------

func createCashfreeOrder(c fiber.Ctx, amountINR float64, req CreateOrderReq) error {
	appID := os.Getenv("CASHFREE_PAYMENT_APP_ID")
	secret := os.Getenv("CASHFREE_PAYMENT_SECRET")

	orderID := fmt.Sprintf("tp_%d", time.Now().UnixMilli())

	phone := strings.TrimPrefix(req.CustomerPhone, "+91")
	phone = strings.TrimSpace(phone)
	if phone == "" {
		phone = "9999999999"
	}
	email := req.CustomerEmail
	if email == "" {
		email = "customer@ticpin.in"
	}
	name := req.CustomerName
	if name == "" {
		name = "Customer"
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"order_id":       orderID,
		"order_amount":   amountINR,
		"order_currency": "INR",
		"customer_details": map[string]interface{}{
			"customer_id":    fmt.Sprintf("cust_%d", time.Now().UnixNano()),
			"customer_name":  name,
			"customer_email": email,
			"customer_phone": phone,
		},
	})

	httpReq, err := http.NewRequest("POST", "https://api.cashfree.com/pg/orders", bytes.NewReader(payload))
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to prepare Cashfree request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-client-id", appID)
	httpReq.Header.Set("x-client-secret", secret)
	httpReq.Header.Set("x-api-version", "2023-08-01")

	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(httpReq)
	if err != nil {
		return utils.ErrorResponse(c, 502, "Cashfree unreachable")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode != 200 {
		msg := "Cashfree order creation failed"
		if m, ok := result["message"].(string); ok {
			msg = m
		}
		return utils.ErrorResponse(c, 502, msg)
	}

	sessionID, _ := result["payment_session_id"].(string)
	cfOrderID, _ := result["order_id"].(string)

	return c.JSON(fiber.Map{
		"status": 200,
		"data": fiber.Map{
			"gateway":            "cashfree",
			"order_id":           cfOrderID,
			"payment_session_id": sessionID,
			"amount":             amountINR,
			"currency":           "INR",
		},
	})
}

// ---- VerifyPayment ----------------------------------------------------------

func VerifyPayment(c fiber.Ctx) error {
	var req VerifyReq
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	switch req.Gateway {
	case "razorpay":
		return verifyRazorpay(c, req)
	case "cashfree":
		return verifyCashfree(c, req)
	default:
		return utils.ErrorResponse(c, 400, "Unknown gateway: "+req.Gateway)
	}
}

func verifyRazorpay(c fiber.Ctx, req VerifyReq) error {
	keySecret := os.Getenv("RAZORPAY_KEY_SECRET")

	// Razorpay signature = HMAC-SHA256(order_id + "|" + payment_id, key_secret)
	message := req.OrderID + "|" + req.PaymentID
	mac := hmac.New(sha256.New, []byte(keySecret))
	mac.Write([]byte(message))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(req.Signature)) {
		return utils.ErrorResponse(c, 400, "Payment verification failed: signature mismatch")
	}

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Payment verified",
		"data":    fiber.Map{"verified": true, "gateway": "razorpay"},
	})
}

func verifyCashfree(c fiber.Ctx, req VerifyReq) error {
	appID := os.Getenv("CASHFREE_PAYMENT_APP_ID")
	secret := os.Getenv("CASHFREE_PAYMENT_SECRET")

	url := fmt.Sprintf("https://api.cashfree.com/pg/orders/%s", req.OrderID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to verify Cashfree payment")
	}
	httpReq.Header.Set("x-client-id", appID)
	httpReq.Header.Set("x-client-secret", secret)
	httpReq.Header.Set("x-api-version", "2023-08-01")

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(httpReq)
	if err != nil {
		return utils.ErrorResponse(c, 502, "Cashfree verification request failed")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	status, _ := result["order_status"].(string)
	if status != "PAID" {
		return utils.ErrorResponse(c, 400, fmt.Sprintf("Payment not completed (status: %s)", status))
	}

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Payment verified",
		"data":    fiber.Map{"verified": true, "gateway": "cashfree"},
	})
}
