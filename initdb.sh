source .env
export JWT_SECRET=$(openssl rand -base64 256)
echo $JWT_SECRET

go run ./entry create_db
go run ./entry create_schema
go run ./entry create_superadmin -e test_super_admin@gmail.com -p password

go run ./entry/main.go
