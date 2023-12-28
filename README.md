# UploadHub - Image Upload and Storage Service

UploadHub is a versatile image upload and storage service developed using Go, Gin, PostgreSQL, RabbitMQ, Minio, and Docker. This project provides user registration and authorization functionalities, allowing users to securely sign up, sign in, upload images, and retrieve images based on their sizes.

## Features

### User Registration and Authorization

Secure user registration and authorization are implemented using JWT Tokens. Users can register, sign in, and obtain JWT Tokens for authenticated access.

### PostgreSQL Integration

User data is persistently stored in a PostgreSQL database, ensuring reliable and durable data storage for user-related information.

### Image Processing

Uploaded images undergo automatic resizing to 75%, 50%, and 25% of their original size. The processed images are then queued (RabbitMQ) for asynchronous handling, enhancing performance and responsiveness.

### Asynchronous Image Storage

Processed images are efficiently stored in Minio, an object storage server. This approach ensures effective management and rapid serving of images while maintaining scalability.

### Dockerized Deployment

The project is containerized with Docker, providing a convenient Docker Compose file for straightforward deployment and scaling. This simplifies the setup process and facilitates easy management of dependencies.

## Routes

### Authentication

- **Sign Up**: `POST /auth/sign-up`
- **Sign In**: `GET /auth/sign-in`
- **Refresh Token**: `GET /auth/refresh`

### User Profile

- **Get User Profile**: `GET /profile/:sub`

### Image Management

- **Upload Image**: `POST /images/upload`
- **Get All Images**: `GET /images/all`
- **Get Images by Size**: `GET /images/by-size`
- **Delete All Images**: `DELETE /images/delete/:name`

- ### Profile
- **Get**: `GET /profile/get/:sub`
  
### Dashboard (Admin Access Only)

- **Get Logs**: `GET /dashboard/logs`
- **Delete Log**: `DELETE /dashboard/logs/:id`

## Usage

1. Clone the repository.
2. Configure the environment variables.
3. Run `docker-compose up` for easy deployment and scaling.

Explore the various routes to leverage the features provided by UploadHub.

Feel free to contribute or report issues on [GitHub](#).

---

*Note: This README assumes you have basic knowledge of Go, Docker, and related technologies.*
