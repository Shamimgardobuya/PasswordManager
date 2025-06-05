# PasswordManager

A simple password manager built with Golang. It supports two versions:
- **CLI version:** Command-line interface for local use.
- **API version:** RESTful API for web or remote access.

## Requirements

- **Go:** v1.24+
- **Database:** MySQL

## Branches

- **CLI version:** [`cli/password_manager`](https://github.com/your-repo/PasswordManager/tree/cli/password_manager)
- **API version:** [`web-app/password-manager`](https://github.com/your-repo/PasswordManager/tree/web-app/password-manager)

## Database Schema

- **Users:**  
    - `id` (Primary Key)  
    - `email`  
    - `password`  

- **Password:**  
    - `id` (Primary Key)  
    - `item` (Site or service name)  
    - `password` (Stored password)  
    - `userId` (Foreign key to Users)

## Setup

1. **Create a `.env` file** in the project root with your database credentials:
        ```
        DB_HOST=your_host
        DB_USER=your_user
        DB_PASSWORD=your_password
        DB_PORT=your_port
        DB_NAME=your_db_name
        ```

2. **Start the application:**  
     Run `password.go` to start the server or CLI.

## API Usage Example

To create a new user via the API, run:

```bash
curl -X POST -H "Content-Type: application/json" \
    -d '{"email":"john@gmail.com","password":"Hello@12347!"}' \
    http://localhost:{port}/users/create/
```

## CLI Usage

For the CLI version, run `password.go` and follow the prompts.

---

Feel free to contribute or open issues for improvements!
