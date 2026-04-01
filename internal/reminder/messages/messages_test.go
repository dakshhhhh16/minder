// SPDX-FileCopyrightText: Copyright 2026 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

package messages

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewEntityReminderMessage_NonNil(t *testing.T) {
	t.Parallel()
	msg, err := NewEntityReminderMessage(uuid.New(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("NewEntityReminderMessage() error = %v", err)
	}
	if msg == nil {
		t.Error("NewEntityReminderMessage() returned nil message")
	}
}

func TestNewEntityReminderMessage_PayloadRoundTrip(t *testing.T) {
	t.Parallel()
	providerID := uuid.New()
	entityID := uuid.New()
	projectID := uuid.New()
	msg, err := NewEntityReminderMessage(providerID, entityID, projectID)
	if err != nil {
		t.Fatalf("NewEntityReminderMessage() error = %v", err)
	}
	evt, err := EntityReminderEventFromMessage(msg)
	if err != nil {
		t.Fatalf("EntityReminderEventFromMessage() error = %v", err)
	}
	if evt.ProviderID != providerID {
		t.Errorf("ProviderID = %v, want %v", evt.ProviderID, providerID)
	}
	if evt.EntityID != entityID {
		t.Errorf("EntityID = %v, want %v", evt.EntityID, entityID)
	}
	if evt.Project != projectID {
		t.Errorf("Project = %v, want %v", evt.Project, projectID)
	}
}

func TestNewEntityReminderMessage_UniqueIDs(t *testing.T) {
	t.Parallel()
	msg1, _ := NewEntityReminderMessage(uuid.New(), uuid.New(), uuid.New())
	msg2, _ := NewEntityReminderMessage(uuid.New(), uuid.New(), uuid.New())
	if msg1.UUID == msg2.UUID {
		t.Error("expected unique message UUIDs")
	}
}

func TestEntityReminderEventFromMessage_InvalidPayload(t *testing.T) {
	t.Parallel()
	msg, _ := NewEntityReminderMessage(uuid.New(), uuid.New(), uuid.New())
	msg.Payload = []byte("not-json")
	_, err := EntityReminderEventFromMessage(msg)
	if err == nil {
		t.Error("EntityReminderEventFromMessage() expected error for invalid payload")
	}
}
