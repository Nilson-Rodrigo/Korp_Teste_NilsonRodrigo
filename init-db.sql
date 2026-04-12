-- Create estoque database
CREATE DATABASE IF NOT EXISTS estoque;

-- Create faturamento database
CREATE DATABASE IF NOT EXISTS faturamento;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE estoque TO postgres;
GRANT ALL PRIVILEGES ON DATABASE faturamento TO postgres;
