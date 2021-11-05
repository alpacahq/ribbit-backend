mkcert -key-file ribbit.com.pem -cert-file ribbit-public.com.pem ribbit.com

# for client request encryption
openssl genrsa -out private_key.pem 1024
openssl rsa -in private_key.pem -outform PEM -pubout -out public_key.pem
