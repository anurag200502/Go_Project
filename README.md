This is a simple Repair Shop Management System implemented in Go with a MongoDB backend. The system allows users to submit, list, update, and delete laptop repair entries.

Install dependencies:
go get -u gofr.dev/pkg/gofr
go get go.mongodb.org/mongo-driver/mongo

Run the application:
go run main.go

Submit a Repair:
curl -X POST -H "Content-Type: application/json" -d '{"customerName": "Anurag Kumar", "repairID": "123456", "deviceName": "Laptop", "status": "Pending"}' http://localhost:8000/repair/submit

List Repairs:
curl http://localhost:8080/repairs/list

Update a Repair:
curl -X PUT -H "Content-Type: application/json" -d '{"status": "Completed"}' http://localhost:8000/repairs/list/123456

Delete a Repair:
curl -X DELETE http://localhost:8080/repair/delete/123456



