package repository

import "time"

type SettlementStatus string

const (
	StatusPending    SettlementStatus = "PENDING"
	StatusProcessed  SettlementStatus = "PROCESSED"
	StatusFailed     SettlementStatus = "FAILED"
	StatusEscrow     SettlementStatus = "ESCROW"
)

type Settlement struct {
	ID             string           `db:"id" json:"id"`
	TelcoReference string           `db:"telco_reference" json:"telco_reference"`
	TxHash         string           `db:"tx_hash" json:"tx_hash"`
	IntentHash     []byte           `db:"intent_hash" json:"-"`
	SenderPhone    string           `db:"sender_phone" json:"sender_phone"`
	ReceiverPhone  string           `db:"receiver_phone" json:"receiver_phone"`
	Amount         float64          `db:"amount" json:"amount"`
	Currency       string           `db:"currency" json:"currency"`
	Status         SettlementStatus `db:"status" json:"status"`
	RawPayload     []byte           `db:"raw_payload" json:"-"`
	MerchantID     string           `db:"merchant_id" json:"merchant_id"`
	Type           string           `db:"type" json:"type"`
	Direction      string           `db:"direction" json:"direction"`
	Counterparty   string           `db:"counterparty" json:"counterparty"`
	FromAddress    string           `db:"from_address" json:"from_address"`
	ToAddress      string           `db:"to_address" json:"to_address"`
	CreatedAt      time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time        `db:"updated_at" json:"updated_at"`
}

type Merchant struct {
	ID            string    `db:"id" json:"id"`
	UserID        string    `db:"user_id" json:"user_id"`
	WalletAddress string    `db:"wallet_address" json:"wallet_address"`
	BusinessName  string    `db:"business_name" json:"business_name"`
	ContactPhone  string    `db:"contact_phone" json:"-"`
	CountryCode   string    `db:"country_code" json:"country"`
	KYCStatus     string    `db:"kyc_status" json:"kyc_status"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

// User — auth identity only. No business/country fields here; use Merchant for that.
type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	Wallet       string    `db:"wallet" json:"wallet"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Node struct {
	ID        string    `db:"id" json:"id"`
	Region    string    `db:"region" json:"region"`
	Status    string    `db:"status" json:"status"`
	LatencyMs int       `db:"latency_ms" json:"latency_ms"`
	Version   string    `db:"version" json:"version"`
	LastSeen  time.Time `db:"last_seen" json:"last_seen"`
}

type TopMerchant struct {
	BusinessName string  `db:"business_name" json:"business_name"`
	Volume24h    float64 `db:"volume_24h" json:"volume_24h"`
	TxCount      int     `db:"tx_count" json:"tx_count"`
	Status       string  `json:"status"`
}