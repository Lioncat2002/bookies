# Bookies
## How to run:
1. Install all dependencies with `go mod download`
2. Create a .env file according to the .example-env
```env
DB_URI=your-postgres-link
API_KEY="supersekret"
TOKEN_LIFESPAN=1
```
3. Migrate the models with `go run migrations/migration.go`
4. Run the app with `go run main.go`
5. The backend should be live at `localhost:8080`
## Endpoints:
[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/22341383-48b44295-ed44-4ec7-8318-bdc05a991852?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D22341383-48b44295-ed44-4ec7-8318-bdc05a991852%26entityType%3Dcollection%26workspaceId%3D8f1d0273-aa92-4236-b8c5-4894eae0a5b7)

## Public deployment
https://spobook.fly.dev/

## ER Diagram:
![er_diagram](https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-Lioncat2002/assets/74904820/4c7f4364-d424-4764-b87c-9f51925bd1bd)
