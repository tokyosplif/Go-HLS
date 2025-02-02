CREATE TABLE IF NOT EXISTS creatives (
         id INT AUTO_INCREMENT PRIMARY KEY,
         campaign_id INT,
         duration INT NOT NULL,
         price DECIMAL(10, 2) NOT NULL,
         playlist_hls VARCHAR(255),
         FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);

INSERT INTO creatives (campaign_id, duration, price, playlist_hls) VALUES
         (1, 20, 10.00, '#EXTINF:5.000, ad1.ts #EXTINF:5.000, ad2.ts #EXTINF:5.000, ad3.ts #EXTINF:5.000, ad4.ts'),
         (1, 20, 12.00, '#EXTINF:5.000, ad5.ts #EXTINF:5.000, ad6.ts #EXTINF:5.000, ad7.ts #EXTINF:5.000, ad8.ts'),
         (1, 25, 9.00, '#EXTINF:5.000, ad9.ts #EXTINF:5.000, ad10.ts #EXTINF:5.000, ad11.ts #EXTINF:5.000, ad12.ts #EXTINF:5.000, ad13.ts'),
         (1, 25, 11.00, '#EXTINF:5.000, ad14.ts #EXTINF:5.000, ad15.ts #EXTINF:5.000, ad16.ts #EXTINF:5.000, ad17.ts #EXTINF:5.000, ad18.ts'),
         (2, 30, 8.50, '#EXTINF:5.000, ad19.ts #EXTINF:5.000, ad20.ts #EXTINF:5.000, ad21.ts #EXTINF:5.000, ad22.ts #EXTINF:5.000, ad23.ts #EXTINF:5.000, ad24.ts'),
         (2, 30, 10.50, '#EXTINF:5.000, ad25.ts #EXTINF:5.000, ad26.ts #EXTINF:5.000, ad27.ts #EXTINF:5.000, ad28.ts #EXTINF:5.000, ad29.ts #EXTINF:5.000, ad30.ts'),
         (2, 25, 7.50, '#EXTINF:5.000, ad31.ts #EXTINF:5.000, ad32.ts #EXTINF:5.000, ad33.ts #EXTINF:5.000, ad34.ts #EXTINF:5.000, ad35.ts'),
         (3, 25, 6.00, '#EXTINF:5.000, ad36.ts #EXTINF:5.000, ad37.ts #EXTINF:5.000, ad38.ts #EXTINF:5.000, ad39.ts #EXTINF:5.000, ad40.ts'),
         (3, 30, 9.50, '#EXTINF:5.000, ad41.ts #EXTINF:5.000, ad42.ts #EXTINF:5.000, ad43.ts #EXTINF:5.000, ad44.ts #EXTINF:5.000, ad45.ts #EXTINF:5.000, ad46.ts'),
         (4, 20, 7.00, '#EXTINF:4.000, ad47.ts #EXTINF:4.000, ad48.ts #EXTINF:4.000, ad49.ts #EXTINF:4.000, ad50.ts'),
         (4, 30, 10.00, '#EXTINF:5.000, ad51.ts #EXTINF:5.000, ad52.ts #EXTINF:5.000, ad53.ts #EXTINF:5.000, ad54.ts #EXTINF:5.000, ad55.ts #EXTINF:5.000, ad56.ts'),
         (5, 25, 7.50, '#EXTINF:5.000, ad57.ts #EXTINF:5.000, ad58.ts #EXTINF:5.000, ad59.ts #EXTINF:5.000, ad60.ts #EXTINF:5.000, ad61.ts'),
         (5, 30, 6.50, '#EXTINF:5.000, ad62.ts #EXTINF:5.000, ad63.ts #EXTINF:5.000, ad64.ts #EXTINF:5.000, ad65.ts #EXTINF:5.000, ad66.ts #EXTINF:5.000, ad67.ts'),
         (6, 20, 6.50, '#EXTINF:5.000, ad68.ts #EXTINF:5.000, ad69.ts #EXTINF:5.000, ad70.ts #EXTINF:5.000, ad71.ts'),
         (6, 25, 8.50, '#EXTINF:5.000, ad72.ts #EXTINF:5.000, ad73.ts #EXTINF:5.000, ad74.ts #EXTINF:5.000, ad75.ts #EXTINF:5.000, ad76.ts'),
         (6, 30, 9.50, '#EXTINF:5.000, ad77.ts #EXTINF:5.000, ad78.ts #EXTINF:5.000, ad79.ts #EXTINF:5.000, ad80.ts #EXTINF:5.000, ad81.ts #EXTINF:5.000, ad82.ts');
