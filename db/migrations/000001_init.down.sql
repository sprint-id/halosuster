DROP INDEX IF EXISTS idx_users_nip;
DROP INDEX IF EXISTS idx_users_id;
DROP INDEX IF EXISTS idx_patients_id;
DROP INDEX IF EXISTS idx_patients_user_id;
DROP INDEX IF EXISTS idx_patients_identity_number;
DROP INDEX IF EXISTS idx_medical_records_id;
DROP INDEX IF EXISTS idx_medical_records_user_id;
DROP INDEX IF EXISTS idx_medical_records_patient_identifier;

DROP TABLE IF EXISTS medical_records;

DROP TABLE IF EXISTS patients;

DROP TABLE IF EXISTS USERS;