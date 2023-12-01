# UploadHub

UploadHub is a simple image upload and storage service with user registration and authorization, built using Go, Gin, MongoDB, RabbitMQ, Minio, and Docker. Users can register, sign in, upload images, and retrieve images based on their sizes.

## Features

- **User Registration and Authorization:** Secure user registration and authorization are implemented using JWT Tokens.

- **MongoDB Integration:** User data is stored in a MongoDB database for persistence.

- **Image Processing:** Uploaded images are automatically resized to 75%, 50%, and 25% of their original size. The processed images are then sent to a queue (RabbitMQ) for asynchronous handling.

- **Asynchronous Image Storage:** Processed images are stored in Minio, an object storage server, to efficiently manage and serve the images.

- **Dockerized Deployment:** The project comes with its own Docker Compose file for easy deployment and scaling.

## Routes

- **Authentication Routes:**
  - `POST /auth/sign-up`: User registration endpoint.
  - `GET /auth/sign-in`: User sign-in endpoint.

- **Image Routes:**
  - `POST /images/upload`: Upload an image.
  - `GET /images/get-all`: Retrieve all uploaded images.
  - `GET /images/get`: Retrieve an image by its size (100, 75, 50, or 25). Requires authentication.

## Code Structure

The project's code is organized as follows:

- **Routes:** Authentication and image-related routes are defined in the respective route groups.

- **Middleware:** Authentication middleware is used to secure image-related routes.

- **Data Models:**
  - `Image`: Represents the structure of an uploaded image.
  - `SignInInput`: Input structure for user sign-in.
  - `SignUpInput`: Input structure for user registration.

## Getting Started

1. Clone the repository: `git clone https://github.com/your-username/UploadHub.git`
2. Navigate to the project directory: `cd UploadHub`
3. Set up environment variables as needed.
4. Run the application using Docker Compose: `docker-compose up -d`

## Dependencies

- [Gin](https://github.com/gin-gonic/gin): Web framework for Go.
- [MongoDB](https://www.mongodb.com/): NoSQL database for user data storage.
- [RabbitMQ](https://www.rabbitmq.com/): Message queue for asynchronous image processing.
- [Minio](https://min.io/): Object storage server for image storage.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This project is licensed under the [MIT License](LICENSE).

