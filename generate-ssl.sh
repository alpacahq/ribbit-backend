mkcert -key-file private_key.pem -cert-file public_cert.pem ribbit.com
openssl rsa -in private_key.pem -out private_key.pem
