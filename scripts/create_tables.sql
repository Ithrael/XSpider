CREATE TABLE websites (
    id INT AUTO_INCREMENT PRIMARY KEY, 
    url VARCHAR(2000) NOT NULL, 
    title VARCHAR(255),
    host	VARCHAR(255) NOT NULL,
    code int,
    finger VARCHAR(255),
    timestamp datetime
);