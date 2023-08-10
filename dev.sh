#!/bin/bash
DB_SOS="postgresql://postgres:6GKePzfr56sBUDC9ctBNq8H36QJcvAsYVqKpiEbF@localhost:5432/sos?schema=public?sslmode=disable"
DB_TEST="user=postgres dbname=sos sslmode=disable password=6GKePzfr56sBUDC9ctBNq8H36QJcvAsYVqKpiEbF"

DATABASE_URL=$DB_TEST go run main.go