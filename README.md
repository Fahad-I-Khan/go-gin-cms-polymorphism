# CMS with Polymorphism using Gin, GORM, and Swagger

This project is a **Content Management System (CMS)** built using Go, leveraging the **Gin web framework**, **GORM ORM**, and **Swagger** for API documentation. It demonstrates the implementation of **polymorphism** to manage articles, videos, tags, and comments.

### Features

- Manage Articles and Videos with associated Tags and Comments.

- Polymorphic relationships implemented with GORM.

- Swagger integration for API documentation.

- Dockerized setup for easy deployment.

## Project Setup

### Clone the Repository

```bash
git clone https://github.com/Fahad-I-Khan/go-gin-cms-polymorphism.git
cd go-gin-cms-polymorphism
```

### Install Dependencies

Ensure all dependencies are installed by reviewing the `import` statements in the `main.go` file and running:

```bash
go mod tidy
```

### Build and Run with Docker

1. Start the database container:

```bash
docker-compose up -d go_db
```

2. Build the application image:

```bash
docker-compose build
```

3. Run the application:

```bash
docker-compose up go-app
```
Swagger UI will be available at: http://localhost:8000/swagger/index.html

### Run Swagger Documentation

If you update the Swagger annotations, regenerate the documentation:

```bash 
swag init
```
**Note**: Swagger may not recognize annotations inside functions. Place them directly above the function.

### API Documentation

#### Articles

**Create Article**

**Endpoint:** `POST /api/v1/articles`

**Request Body**:

```json
{
  "title": "Article Title",
  "content": "This is the content of the article.",
  "tags": [
    { "name": "Tag1" },
    { "name": "Tag2" }
  ],
  "comments": [
    { "content": "This is a comment." }
  ]
}
```
**Update Article**

**Endpoint:** `PUT /api/v1/articles/{id}`

**Request Body**:

```json
{
  "title": "Updated Title",
  "content": "Updated content.",
  "tags": [
    { "name": "Updated Tag" }
  ],
  "comments": [
    { "content": "Updated comment." }
  ]
}
```
**Get All Articles**

**Endpoint**: `GET /api/v1/articles`

#### Videos

**Create Video**

**Endpoint**: `POST /api/v1/videos`

**Request Body**:

```json
{
  "title": "Video Title",
  "url": "http://example.com/video.mp4",
  "tags": [
    { "name": "Tag1" },
    { "name": "Tag2" }
  ],
  "comments": [
    { "content": "This is a comment." }
  ]
}
```

**Update Video**

**Endpoint**: `PUT /api/v1/videos/{id}`

**Request Body**:

```json
{
  "title": "Updated Video Title",
  "url": "http://example.com/updated-video.mp4",
  "tags": [
    { "name": "Updated Tag" }
  ],
  "comments": [
    { "content": "Updated comment." }
  ]
}
```

**Get All Videos**

**Endpoint**: `GET /api/v1/videos`

### Polymorphism in GORM

Polymorphism in this project allows Tags and Comments to be associated with multiple models (Articles and Videos).

#### Implementation

- **Each polymorphic model (e.g.,** `Tag` **and** `Comment`**) includes:**
  - A `TaggableID` or `CommentableID `to store the ID of the associated record. The `gorm:"index"` annotation ensures the fields are indexed in the database, improving query performance when filtering by `TaggableID` or `CommentableID`.
  - A `TaggableType` or `CommentableType` to store the type of the associated model (e.g., `articles` or `videos`).
- In GORM, use the `gorm:"polymorphic:<field_name>;"` tag for relationships. This tag tells GORM to use the `<field_name>ID` and `<field_name>Type` fields to establish the polymorphic association.

#### Example

**Tag Model**:

```go 
type Tag struct {
  BaseModel
  Name         string `json:"name"`
  TaggableID   uint   `gorm:"index"`
  TaggableType string `gorm:"index"`
}
```
**Article Model**:

```go
type Article struct {
  BaseModel
  Title    string    `json:"title"`
  Content  string    `json:"content"`
  Tags     []Tag     `gorm:"polymorphic:Taggable;"`
  Comments []Comment `gorm:"polymorphic:Commentable;"`
}
```
The `gorm:"polymorphic:Taggable;" `tag in the `Article` model links it with the `Tag` model, ensuring that each tag is correctly associated with its parent record. Similar logic applies to comments.

GORM automatically handles these relationships, enabling queries like:

```go
db.Preload("Tags").Preload("Comments").Find(&articles)
```

### Why Use /api/v1/ in Endpoints

Using `/api/v1/` in your endpoints helps version your API. This ensures backward compatibility for existing clients when you release updates or make breaking changes in future versions of the API. It is a widely accepted industry practice for maintainable and scalable API development.

### Conclusion

This project demonstrates the power of Go, Gin, and GORM for building robust APIs. The polymorphic relationships make it easy to manage shared entities across different models. Swagger adds a user-friendly layer for exploring the API.