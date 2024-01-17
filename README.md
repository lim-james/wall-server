# Wall Server


# Wall Posts API

## Base URL

`https://ec2-52-77-251-29.ap-southeast-1.compute.amazonaws.com`

# API Documentation

## Overview

This document outlines the RESTful API endpoints for [Your Server Name]. The server is built using the Gin framework in Go and interacts with a MySQL database.

## Base URL

All API endpoints are relative to the base URL: `/api`

---

## Authentication Endpoints

### Sign Up
- **Path**: `/api/u/signup`
- **Method**: POST
- **Body**: 
  ```json
  {
      "username": "user123",
      "password_hash": "hashedpassword"
  }
  ```
- **Example Request** (cURL):
  ```
  curl -X POST https://ec2-52-77-251-29.ap-southeast-1.compute.amazonaws.com/api/u/signup \
  -H 'Content-Type: application/json' \
  -d '{"username": "user123", "password_hash": "hashedpassword"}'
  ```
- **Response**:
  - Success: Status code 200 with user details
  - Error: Status code 4xx/5xx with error message

### Login
- **Path**: `/api/u/login`
- **Method**: POST
- **Body**: 
  ```json
  {
      "username": "user123",
      "password_hash": "hashedpassword"
  }
  ```
- **Example Request** (cURL):
  ```
  curl -X POST https://ec2-52-77-251-29.ap-southeast-1.compute.amazonaws.com/api/u/login \
  -H 'Content-Type: application/json' \
  -d '{"username": "user123", "password_hash": "hashedpassword"}'
  ```
- **Response**:
  - Success: Status code 200 with JWT token
  - Error: Status code 401 for unauthorized access

---

## Post Endpoints

### Create Post
- **Path**: `/api/p/`
- **Method**: POST
- **Auth Required**: Yes
- **Body**:
  ```json
  {
      "title": "Sample Title",
      "body": "This is a sample post."
  }
  ```
- **Example Request** (cURL):
  ```
  curl -X POST https://ec2-52-77-251-29.ap-southeast-1.compute.amazonaws.com/api/p/ \
  -H 'Authorization: Bearer [JWT Token]' \
  -H 'Content-Type: application/json' \
  -d '{"title": "Sample Title", "body": "This is a sample post."}'
  ```
- **Response**:
  - Success: Status code 200 with created post details
  - Error: Status code 4xx/5xx with error message

### Like Post
- **Path**: `/api/p/:post_id/like`
- **Method**: POST
- **Auth Required**: Yes
- **Example Request** (cURL):
  ```
  curl -X POST https://ec2-52-77-251-29.ap-southeast-1.compute.amazonaws.com/api/p/123/like \
  -H 'Authorization: Bearer [JWT Token]'
  ```
  (where 123 is the post_id)
- **Response**:
  - Success: Status code 200 with updated like status
  - Error: Status code 4xx/5xx with error message

### Unlike Post
- **Path**: `/api/p/:post_id/unlike`
- **Method**: POST
- **Auth Required**: Yes
- **Example Request** (cURL):
  ```
  curl -X POST https://ec2-52-77-251-29.ap-southeast-1.compute.amazonaws.com`/api/p/123/unlike \
  -H 'Authorization: Bearer [JWT Token]'
  ```
  (where 123 is the post_id)
- **Response**:
  - Success: Status code 200 with updated like status
  - Error: Status code 4xx/5xx with error message

---

## Error Codes

- `4xx`: Client errors (e.g., bad request, unauthorized, not found)
- `5xx`: Server errors (e.g., internal server error)
