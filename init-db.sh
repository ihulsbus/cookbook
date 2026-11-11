#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER recipe WITH PASSWORD 'recipe_pass';
    CREATE USER ingredient WITH PASSWORD 'ingredient_pass';
    CREATE USER instruction WITH PASSWORD 'instruction_pass';
    CREATE USER metadata WITH PASSWORD 'metadata_pass';
    CREATE USER image WITH PASSWORD 'image_pass';
    CREATE USER search WITH PASSWORD 'search_pass';
    -- "users", not "user": USER is a reserved word in Postgres.
    CREATE USER users WITH PASSWORD 'users_pass';

    CREATE DATABASE recipe OWNER recipe;
    CREATE DATABASE ingredient OWNER ingredient;
    CREATE DATABASE instruction OWNER instruction;
    CREATE DATABASE metadata OWNER metadata;
    CREATE DATABASE image OWNER image;
    CREATE DATABASE search OWNER search;
    CREATE DATABASE users OWNER users;

    GRANT ALL PRIVILEGES ON DATABASE recipe TO recipe;
    GRANT ALL PRIVILEGES ON DATABASE ingredient TO ingredient;
    GRANT ALL PRIVILEGES ON DATABASE instruction TO instruction;
    GRANT ALL PRIVILEGES ON DATABASE metadata TO metadata;
    GRANT ALL PRIVILEGES ON DATABASE image TO image;
    GRANT ALL PRIVILEGES ON DATABASE search TO search;
    GRANT ALL PRIVILEGES ON DATABASE users TO users;
EOSQL

echo "All databases and users created successfully!"
