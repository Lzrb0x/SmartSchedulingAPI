package domain

import "time"

type Tenant struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UserRole string

const (
	RoleOwner    UserRole = "owner"
	RoleClient   UserRole = "client"
	RoleEmployee UserRole = "employee"
)

type User struct {
	ID        int64     `db:"id" json:"id"`
	TenantID  int64     `db:"tenant_id" json:"tenant_id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password_hash" json:"-"`
	Role      UserRole  `db:"role" json:"role"`
	Active    bool      `db:"active" json:"active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Barbershop struct {
	ID              int64     `db:"id" json:"id"`
	TenantID        int64     `db:"tenant_id" json:"tenant_id"`
	Name            string    `db:"name" json:"name"`
	Description     string    `db:"description" json:"description"`
	Characteristics string    `db:"characteristics" json:"characteristics"`
	Address         string    `db:"address" json:"address"`
	Contact         string    `db:"contact" json:"contact"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

type Service struct {
	ID           int64     `db:"id" json:"id"`
	TenantID     int64     `db:"tenant_id" json:"tenant_id"`
	BarbershopID int64     `db:"barbershop_id" json:"barbershop_id"`
	Name         string    `db:"name" json:"name"`
	Description  string    `db:"description" json:"description"`
	Price        int64     `db:"price_cents" json:"price_cents"`
	DurationMin  int       `db:"duration_min" json:"duration_min"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Professional struct {
	ID           int64     `db:"id" json:"id"`
	TenantID     int64     `db:"tenant_id" json:"tenant_id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	BarbershopID int64     `db:"barbershop_id" json:"barbershop_id"`
	Specialties  string    `db:"specialties" json:"specialties"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Agenda struct {
	ID             int64     `db:"id" json:"id"`
	ProfessionalID int64     `db:"professional_id" json:"professional_id"`
	TenantID       int64     `db:"tenant_id" json:"tenant_id"`
	Date           time.Time `db:"date" json:"date"`
	StartTime      time.Time `db:"start_time" json:"start_time"`
	EndTime        time.Time `db:"end_time" json:"end_time"`
	Type           string    `db:"type" json:"type"`
}

type Booking struct {
	ID             int64     `db:"id" json:"id"`
	TenantID       int64     `db:"tenant_id" json:"tenant_id"`
	ClientID       int64     `db:"client_id" json:"client_id"`
	ProfessionalID int64     `db:"professional_id" json:"professional_id"`
	ServiceID      int64     `db:"service_id" json:"service_id"`
	StartTime      time.Time `db:"start_time" json:"start_time"`
	EndTime        time.Time `db:"end_time" json:"end_time"`
	Status         string    `db:"status" json:"status"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}
