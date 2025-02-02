CREATE TABLE IF NOT EXISTS source_campaigns (
        source_id INT,
        campaign_id INT,
        PRIMARY KEY (source_id, campaign_id),
        FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE,
        FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);

INSERT INTO source_campaigns (source_id, campaign_id) VALUES
          (1, 1),
          (1, 2),
          (2, 3),
          (2, 4),
          (3, 5),
          (3, 6);
