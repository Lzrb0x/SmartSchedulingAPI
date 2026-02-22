-- +goose Up
CREATE TABLE tenants (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('owner','client','employee')),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, email)
);

CREATE TABLE barbershops (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    characteristics TEXT,
    address TEXT,
    contact TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, name)
);

CREATE TABLE services (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    barbershop_id BIGINT NOT NULL REFERENCES barbershops(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    price_cents BIGINT NOT NULL,
    duration_min INT NOT NULL CHECK (duration_min > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, barbershop_id, name)
);

CREATE TABLE professionals (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    barbershop_id BIGINT NOT NULL REFERENCES barbershops(id) ON DELETE CASCADE,
    specialties TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, user_id)
);

CREATE TABLE agendas (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    professional_id BIGINT NOT NULL REFERENCES professionals(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('work','break','off')),
    CHECK (start_time < end_time)
);

CREATE TABLE bookings (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id BIGINT NOT NULL REFERENCES users(id),
    professional_id BIGINT NOT NULL REFERENCES professionals(id),
    service_id BIGINT NOT NULL REFERENCES services(id),
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'confirmed' CHECK (status IN ('pending','confirmed','canceled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (start_time < end_time)
);

CREATE UNIQUE INDEX bookings_professional_time_unique
    ON bookings (professional_id, start_time, end_time);

CREATE INDEX bookings_tenant_idx ON bookings (tenant_id);
CREATE INDEX agendas_professional_idx ON agendas (professional_id, start_time);

-- +goose Down
DROP TABLE bookings;
DROP TABLE agendas;
DROP TABLE professionals;
DROP TABLE services;
DROP TABLE barbershops;
DROP TABLE users;
DROP TABLE tenants;
