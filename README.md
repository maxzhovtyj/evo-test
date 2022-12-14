# EVO TEST TASK

---
### Maksym Zhovtaniuk (@maxzhovtyj)

## Steps to run the application

* Run docker-compose command in root app folder:
```
docker compose up --build server 
```

* Open a new terminal and move to migrator directory
```shell
cd ./migrator
```

* Build migrator image
```
docker build -t evo-migrator .
```

* Apply migrations to database
```
docker run --network host evo-migrator -path=/schema -database "postgresql://postgres:postgres@localhost:5555/postgres?sslmode=disable" up
```

---

## Usage
##### To test the app use either [test-api file](test-api.http) or [swagger docs](#swagger-docs)

---

# swagger-docs
##### http://localhost:8089/swagger/index.html

![swag-img](swagger-example.png)