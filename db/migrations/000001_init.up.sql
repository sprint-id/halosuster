BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS USERS (
    id UUID PRIMARY KEY,
    nip VARCHAR UNIQUE,
    name VARCHAR,
    password VARCHAR,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES USERS(id) ON DELETE CASCADE,
    identity_number VARCHAR UNIQUE NOT NULL,
    phone_number VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    birth_date VARCHAR NOT NULL,
    gender VARCHAR NOT NULL,
    identity_card_scan_img VARCHAR NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS medical_records (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES USERS(id) ON DELETE CASCADE,
    patient_identifier VARCHAR NOT NULL REFERENCES patients(identity_number) ON DELETE CASCADE,
    symptoms VARCHAR NOT NULL,
    medications VARCHAR NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE INDEX idx_users_nip ON users (nip);
CREATE INDEX idx_users_id ON users (id);
CREATE INDEX idx_patients_id ON patients (id);
CREATE INDEX idx_patients_user_id ON patients (user_id);
CREATE INDEX idx_patients_identity_number ON patients (identity_number);
CREATE INDEX idx_medical_records_id ON medical_records (id);
CREATE INDEX idx_medical_records_user_id ON medical_records (user_id);
CREATE INDEX idx_medical_records_patient_identifier ON medical_records (patient_identifier);

COMMIT TRANSACTION;