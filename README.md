# u-Task Backend

## How to Run

```bash
go run main.go
```

## Make a Request From Postman with Environment

- Development Environment
  <br>
  <b>localhost:3000/</b>
- Production Environment
  <b>https://utask-backend-production.up.railway.app/</b>

## Environment Variables Example

[Environment Variables](config.env.example)

## Endpoint

### Authentication

- <b>/auth/login-gl</b>
  <br>
  Endpoint ini digunakan untuk login
- <b>/auth/sign-in</b>
  <br>
  Endpoint ini digunakan untuk sign-in

### Calendar

- <b>/calendar/</b>
  <br>
  Endpoint ini digunakan untuk melihat semua event di Google Calendar.
- <b>/calendar/:eventId</b>
  <br>
  Endpoint ini digunakan untuk melihat detail dari satu event.
- <b>/calendar/message/:eventId</b>
  <br>
  Endpoint ini digunakan untuk mengirimkan notification berupa SMS Message
