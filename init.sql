CREATE TABLE IF NOT EXISTS Users (
                                     ID INT AUTO_INCREMENT PRIMARY KEY,
                                     Login VARCHAR(255) NOT NULL,
                                     Password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Tasks (
                                     task_id INT AUTO_INCREMENT PRIMARY KEY,
                                     expression VARCHAR(255),
                                     status VARCHAR(255),
                                     answer VARCHAR(255),
                                     login VARCHAR(255)
);
