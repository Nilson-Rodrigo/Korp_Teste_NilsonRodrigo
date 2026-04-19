-- Create estoque database
CREATE DATABASE estoque;

-- Create faturamento database  
CREATE DATABASE faturamento;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE estoque TO postgres;
GRANT ALL PRIVILEGES ON DATABASE faturamento TO postgres;
