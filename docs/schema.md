# Database Schema

This document outlines the database schema for the application.

## Tables

### 1. `users`

Stores information about human actors who can log in.

- `id` (UUID, Primary Key, Not Null)
- `username` (VARCHAR, Unique, Not Null)
- `email` (VARCHAR, Unique, Not Null)
- `password_hash` (VARCHAR, Not Null)
- `created_at` (TIMESTAMPTZ, Not Null, Default `NOW()`)
- `updated_at` (TIMESTAMPTZ, Not Null, Default `NOW()`)

## Notes

- All primary keys are UUIDs.
- Timestamps (`created_at`, `updated_at`) are stored with time zone information (`TIMESTAMPTZ`).
- Indexes will be added for foreign keys and frequently queried columns (e.g., `users.username`, `users.email`).
