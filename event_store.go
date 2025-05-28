package accounting

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EventType constants for different event types
const (
	EventCreateAccount      = "CREATE_ACCOUNT"
	EventUpdateAccount      = "UPDATE_ACCOUNT"
	EventCreateTransaction  = "CREATE_TRANSACTION"
	EventPostTransaction    = "POST_TRANSACTION"
	EventReverseTransaction = "REVERSE_TRANSACTION"
	EventCreatePeriod       = "CREATE_PERIOD"
	EventClosePeriod        = "CLOSE_PERIOD"
	EventReconcile          = "RECONCILE"
)

// EventStore manages the append-only event log
type EventStore struct {
	storage *Storage
}

// NewEventStore creates a new event store
func NewEventStore(storage *Storage) *EventStore {
	return &EventStore{storage: storage}
}

// CreateEvent creates a new journal event
func (es *EventStore) CreateEvent(eventType string, payload interface{}, validTime time.Time, userID string) (*JournalEvent, error) {
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	event := &JournalEvent{
		ID:              uuid.New().String(),
		EventType:       eventType,
		Payload:         payloadData,
		ValidTime:       validTime,
		TransactionTime: time.Now(),
		UserID:          userID,
	}

	if err := es.storage.AppendEvent(event); err != nil {
		return nil, fmt.Errorf("failed to append event: %w", err)
	}

	return event, nil
}

// GetEvents retrieves events within a time range
func (es *EventStore) GetEvents(from, to time.Time) ([]*JournalEvent, error) {
	return es.storage.GetEvents(from, to)
}

// ReplayEvents replays events to rebuild state
func (es *EventStore) ReplayEvents(from, to time.Time, handler func(*JournalEvent) error) error {
	events, err := es.GetEvents(from, to)
	if err != nil {
		return fmt.Errorf("failed to get events: %w", err)
	}

	for _, event := range events {
		if err := handler(event); err != nil {
			return fmt.Errorf("failed to handle event %s: %w", event.ID, err)
		}
	}

	return nil
}

// AccountCreatedEvent represents an account creation event payload
type AccountCreatedEvent struct {
	Account *Account `json:"account"`
}

// TransactionCreatedEvent represents a transaction creation event payload
type TransactionCreatedEvent struct {
	Transaction *Transaction `json:"transaction"`
}

// TransactionPostedEvent represents a transaction posting event payload
type TransactionPostedEvent struct {
	TransactionID string    `json:"transaction_id"`
	PostedAt      time.Time `json:"posted_at"`
	Entries       []Entry   `json:"entries"`
}

// EventProcessor processes events and updates projections
type EventProcessor struct {
	storage *Storage
}

// NewEventProcessor creates a new event processor
func NewEventProcessor(storage *Storage) *EventProcessor {
	return &EventProcessor{storage: storage}
}

// ProcessEvent processes a single event and updates relevant projections
func (ep *EventProcessor) ProcessEvent(event *JournalEvent) error {
	switch event.EventType {
	case EventCreateAccount:
		return ep.handleAccountCreated(event)
	case EventCreateTransaction:
		return ep.handleTransactionCreated(event)
	case EventPostTransaction:
		return ep.handleTransactionPosted(event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}

func (ep *EventProcessor) handleAccountCreated(event *JournalEvent) error {
	var payload AccountCreatedEvent
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal account created event: %w", err)
	}

	return ep.storage.SaveAccount(payload.Account)
}

func (ep *EventProcessor) handleTransactionCreated(event *JournalEvent) error {
	var payload TransactionCreatedEvent
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal transaction created event: %w", err)
	}

	return ep.storage.SaveTransaction(payload.Transaction)
}

func (ep *EventProcessor) handleTransactionPosted(event *JournalEvent) error {
	var payload TransactionPostedEvent
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal transaction posted event: %w", err)
	}

	// Update transaction status
	txn, err := ep.storage.GetTransaction(payload.TransactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	txn.Status = Posted
	txn.UpdatedAt = time.Now()

	if err := ep.storage.SaveTransaction(txn); err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// Save all entries
	for _, entry := range payload.Entries {
		if err := ep.storage.SaveEntry(&entry); err != nil {
			return fmt.Errorf("failed to save entry: %w", err)
		}
	}

	return nil
}
