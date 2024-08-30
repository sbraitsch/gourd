#!/bin/bash

# Set environment variables for the PostgreSQL container
DB_NAME="gourd_db"
DB_USER="local"
DB_PASSWORD="pwd"
DB_PORT="5432"

# Run PostgreSQL Docker container
docker run -d \
  --name gourd_db \
  -e POSTGRES_DB=$DB_NAME \
  -e POSTGRES_USER=$DB_USER \
  -e POSTGRES_PASSWORD=$DB_PASSWORD \
  -p $DB_PORT:5432 \
  postgres