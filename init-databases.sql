-- Create databases for microservices
CREATE DATABASE wallet;
CREATE DATABASE transaction;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE wallet TO postgres;
GRANT ALL PRIVILEGES ON DATABASE transaction TO postgres;