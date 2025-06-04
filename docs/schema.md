# Database Schema

This document outlines the database schema for the application.

## Tables

### 1. `users`

Stores information about human actors who can log in.

-   `id` (UUID, Primary Key, Not Null)
-   `username` (VARCHAR, Unique, Not Null)
-   `email` (VARCHAR, Unique, Not Null)
-   `password_hash` (VARCHAR, Not Null)
-   `created_at` (TIMESTAMPTZ, Not Null, Default `NOW()`)
-   `updated_at` (TIMESTAMPTZ, Not Null, Default `NOW()`)

### 2. `vendors`

Stores information about vendor entities, managed by users.

-   `id` (UUID, Primary Key, Not Null)
-   `name` (VARCHAR, Not Null)
-   `description` (TEXT, Nullable)
-   `user_id` (UUID, Foreign Key referencing `users.id`, Not Null)
-   `created_at` (TIMESTAMPTZ, Not Null, Default `NOW()`)
-   `updated_at` (TIMESTAMPTZ, Not Null, Default `NOW()`)

### 3. `merchants`

Stores information about merchant entities, managed by users.

-   `id` (UUID, Primary Key, Not Null)
-   `name` (VARCHAR, Not Null)
-   `description` (TEXT, Nullable)
-   `user_id` (UUID, Foreign Key referencing `users.id`, Not Null)
-   `created_at` (TIMESTAMPTZ, Not Null, Default `NOW()`)
-   `updated_at` (TIMESTAMPTZ, Not Null, Default `NOW()`)

## Notes

-   All primary keys are UUIDs.
-   Timestamps (`created_at`, `updated_at`) are stored with time zone information (`TIMESTAMPTZ`).
-   `user_id` in `vendors` and `merchants` tables establishes ownership by a user.
-   Indexes will be added for foreign keys and frequently queried columns (e.g., `users.username`, `users.email`).
