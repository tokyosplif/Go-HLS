CREATE TABLE IF NOT EXISTS campaigns (
      id INT AUTO_INCREMENT PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      start_time DATETIME NOT NULL,
      end_time DATETIME NOT NULL,
      source_id INT,
      FOREIGN KEY (source_id) REFERENCES sources(id)
);

INSERT INTO campaigns (id, name, start_time, end_time, source_id) VALUES
   (1, 'Campaign 1', '2025-01-01 00:00:00', '2025-12-31 23:59:59', 1),
   (2, 'Campaign 2', '2025-01-01 00:00:00', '2025-06-30 23:59:59', 1),
   (3, 'Campaign 3', '2025-02-01 00:00:00', '2025-12-31 23:59:59', 2),
   (4, 'Campaign 4', '2025-01-01 00:00:00', '2025-06-30 23:59:59', 2),
   (5, 'Campaign 5', '2025-01-01 00:00:00', '2025-12-31 23:59:59', 3),
   (6, 'Campaign 6', '2025-01-01 00:00:00', '2025-07-31 23:59:59', 3);