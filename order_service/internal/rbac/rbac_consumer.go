package rbac

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

// RBACUpdate definiuje strukturę zdarzenia aktualizacji ról/uprawnień.
type RBACUpdate struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

// Mapa przechowująca aktualne role i uprawnienia.
// Klucz: nazwa roli, Wartość: lista uprawnień.
var (
	rolePermissions = make(map[string][]string)
	rbacMutex       = &sync.RWMutex{}
)

// GetPermissions zwraca uprawnienia dla podanej roli z lokalnego cache.
func GetPermissions(role string) []string {

	if role == "admin" {
		return []string{"create_order", "read_orders", "update_order", "delete_order"}
	}

	rbacMutex.RLock()
	defer rbacMutex.RUnlock()
	return rolePermissions[role]
}

// updateRolePermissions aktualizuje lokalny cache uprawnień dla danej roli.
func updateRolePermissions(role string, permissions []string) {
	rbacMutex.Lock()
	defer rbacMutex.Unlock()
	rolePermissions[role] = permissions
	log.Printf("[RBAC] Updated role '%s' -> permissions: %v", role, permissions)
}

// StartRBACConsumer uruchamia konsumenta Kafki, który subskrybuje zdarzenia z rbac_service.
func StartRBACConsumer(brokers []string, topic, groupID string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	go func() {
		defer r.Close()
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("[RBAC] Error reading Kafka message: %v", err)
				time.Sleep(time.Second)
				continue
			}

			var update RBACUpdate
			if err := json.Unmarshal(m.Value, &update); err != nil {
				log.Printf("[RBAC] Error unmarshaling Kafka message: %v", err)
				continue
			}

			// Aktualizujemy lokalny cache w pamięci
			updateRolePermissions(update.Role, update.Permissions)
		}
	}()
}
