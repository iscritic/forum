-- migration.sql

INSERT INTO users (username, email, password, role) VALUES
    ('alice', 'alice@example.com', 'hashed_password_1', 'user'),
    ('bob', 'bob@example.com', 'hashed_password_2', 'user'),
    ('carol', 'carol@example.com', 'hashed_password_3', 'admin');

INSERT INTO posts (title, content, author_id, category_id) VALUES
    ('The Rise of AI', 'Exploring the advancements in artificial intelligence.', 1, 1),
    ('Healthy Eating Tips', 'A guide to maintaining a balanced diet.', 2, 4),
    ('Top Destinations for 2024', 'A list of must-visit places for the upcoming year.', 3, 3),
    ('The Future of Quantum Computing', 'Understanding the potential of quantum computing technology.', 1, 1),
    ('Benefits of Regular Exercise', 'Why staying active is crucial for your health.', 2, 5);

-- Comments
INSERT INTO comments (post_id, content, author_id) VALUES
    (1, 'Very informative article!', 2),
    (2, 'I found these tips really helpful.', 3),
    (3, 'Looking forward to visiting some of these places.', 1),
    (4, 'Excited about the future of tech!', 3),
    (5, 'Great article on fitness.', 1);

-- Likes
INSERT INTO likes (user_id, post_id) VALUES
    (1, 1),
    (2, 2),
    (3, 3),
    (1, 4),
    (2, 5);

-- Dislikes
INSERT INTO dislikes (user_id, post_id) VALUES
    (2, 1),
    (3, 2),
    (1, 3);
