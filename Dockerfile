Fversion: '3.8'

 services:
   postgres:
     image: postgres:latest
     environment:
       POSTGRES_USER: postgres
       POSTGRES_PASSWORD: secret
       POSTGRES_DB: orders
     ports:
       - "5432:5432"
     volumes:
       - postgres_data:/var/lib/postgresql/data
     networks:
       - orders-network

   order-service:
     build: .
     ports:
       - "3000:3000"
     environment:
       JWT_SECRET: "your_jwt_secret_123"
       DB_HOST: "postgres"
       DB_USER: "postgres"
       DB_PASSWORD: "secret"
       DB_NAME: "orders"
     depends_on:
       - postgres
     networks:
       - orders-network

 volumes:
   postgres_data:

 networks:
   orders-network:
     driver: bridge