CREATE TABLE class (
    ClassID SERIAL PRIMARY KEY,
    TeacherID INT,
    Name VARCHAR(255),
    Description TEXT,
    Created_at DATE
);