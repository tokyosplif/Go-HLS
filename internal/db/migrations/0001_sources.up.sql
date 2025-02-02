CREATE TABLE IF NOT EXISTS sources (
     id INT AUTO_INCREMENT PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
     status ENUM('active', 'inactive') NOT NULL
);

INSERT INTO sources (name, status) VALUES
        ('Source 1', 'active'),
        ('Source 2', 'inactive'),
        ('Source 3', 'active');
