
INSERT INTO category (name)
VALUES
    ('Off-Topic'),
    ('Movies'),
    ('News'),
    ('Pets'),
    ('Games')
ON CONFLICT(name) DO NOTHING;

INSERT INTO users (username, email, password, role) VALUES
    ('Dias', 'diass@gmail.com', '88888888', 'user'),
    ('Arthur', 'iscritic@gmail.com', '87654321', 'user'),
    ('Anelkhan', 'goddamness@gmail.com', '12345678', 'user');

INSERT INTO posts (title, content, author_id, category_id) VALUES
    ('The Gopher, the God, and the Joker Girl: A Witcher Tale', 'In the mystical land of Tashkent, Mufasa, a renowned Witcher of the Gopher School, finds himself entwined in a tale of destiny and magic. Alongside his best friend, Joekerr Girl, a bard known for her wit and mischief, and SultanBek, the radiant God of the Dawn, they form an unexpected trio. Mufasa, with his unassuming strength and wisdom, and SultanBek, with his divine charisma and fondness for Greek mythology, soon become close friends. Their adventures are filled with laughter, peril, and the occasional clash of personalities. Joekerr Girl, despite her playful nature, harbors a deep-seated envy towards the bond forming between Mufasa and SultanBek. When Mufasa gently rejects her romantic advances, Joekerr Girl, hurt and vengeful, succumbs to the dark allure of the Shadow. Transformed into a formidable antagonist, her once-jovial demeanor becomes a harbinger of chaos and betrayal. The once-laughing trio is now fractured, with Joekerr Girl unleashing turmoil upon her former friends. As the Gopher and the God prepare for a final confrontation, they must grapple with their own emotions and the darkness that has corrupted their friend. The stage is set for an epic showdown that will test the limits of their strength and the bonds of their friendship, as they face the darkness together in a dramatic twist of fate.', 1, 5),
    ('To All Those We Have Lost', 'This message goes out to the heroes, the friends, the loved ones that have fallen in our battles. May their memories be a beacon of hope and a reminder of the strength we hold within. Their sacrifices will not be forgotten, their legacy will endure. May their souls find peace, and their spirits guide us in the darkness. Though they are gone, their stories live on in our hearts, a testament to the bonds of friendship and the courage that shines even in the face of despair. We honor their memory by fighting for the light, by continuing their fight for a better world, and by carrying their torch until the very end. Farewell, dear friends, your light will forever guide us.', 2, 3),
    ('love story of witcher🌚', 'Mufasa of Tashkent, a Witcher of the Gopher School, preferred his titanium boomerang to silver swords. He was known for his practicality, his yellow tracksuit, and his stoicism. But then SultanBek, the shimmering golden God of the Dawn, entered his life. SultanBek, with his love for Greek mythology, flamboyant entrances, and an infectious laugh, was a whirlwind of energy. He charmed Mufasa with his theatrical tales, his genuine interest in the Witcher’s stories, and his warm, golden gaze. Their adventures, filled with laughter and monsters, became a symphony of contrasting personalities. Mufasa, the stoic, found himself drawn to SultanBek’s warmth and unbridled enthusiasm, while SultanBek, the God of Dawn, was captivated by the Witcher’s quiet strength and unwavering resolve. One starry night, after vanquishing a particularly nasty pack of werewolves, they found themselves lying on a grassy hill, gazing at the Milky Way. The silence, broken only by the chirping of crickets, felt charged with unspoken emotions. "You are unlike any other I have ever met, Witcher," SultanBek whispered, his golden glow illuminating Mufasa’s face. "Your strength, your quiet humor, it captivates me." Mufasa, surprised by the declaration, found himself captivated by the God’s vulnerability. "And I, God," he replied, his voice rough with emotion, "find myself charmed by your… unique perspective on life." Their eyes met, and the moment stretched, filled with unspoken promises. Finally, SultanBek leaned in, his golden lips meeting Mufasa’s in a kiss that was both tender and powerful. The Gopher and the God, their love story an unexpected symphony of contrasts, had found solace in each other’s company. They had found a love that transcended the boundaries of reality, a love that was as unique and powerful as the clash of their personalities.🌚🌚🌚', 3, 1),
    ('The Little Dragon Who Loved Rainbows', 'In a hidden valley where waterfalls cascaded and rainbows danced, lived a tiny dragon named Sparky. Unlike his fire-breathing kin, Sparky adored the colors of the rainbow. He would chase them across the sky, his scales shimmering with a kaleidoscope of hues. But the other dragons teased him, calling him "Rainbow Sparkles." Sparky felt lonely, until one day, he met a young girl named Lily, who loved rainbows just as much. Together, they explored the valley, chasing rainbows and sharing stories. Sparky learned that being different was beautiful, and Lily learned that even dragons could have hearts full of color.', 1, 4),
    ('Mads Mikkelsen: A Chameleon of Talent', 'Mads Mikkelsen is a master of transformation.  He effortlessly slips into diverse roles, from the chilling Hannibal Lecter to the charming Tristan in "The Kings Speech.  His intensity and captivating gaze draw us into each character, leaving us both enthralled and unnerved. Mikkelsens versatility proves that true talent transcends genre and form, making him a force to be reckoned with in the world of film.', 2, 2);

INSERT INTO comments (post_id, content, author_id) VALUES
    (1, 'Wow, this is such a creative Witcher tale! I love the unexpected trio and the twist with Joekerr Girl.', 2),
    (1, 'This sounds like it could be the start of a really epic fantasy series. I need more!', 3),
    (2, 'This is such a powerful tribute. Its important to remember those who have fought for us.', 1),
    (2, 'Beautifully written. This message will resonate with so many people.', 3),
    (3, 'This is such a heartwarming love story! I never would have guessed a Gopher and a God would fall in love!', 1),
    (3, 'This is a story that could only happen in a world filled with magic and wonder. I love it!🌈🌈🌈', 2),
    (4, 'This little dragon is adorable! I love how he embraces his uniqueness.', 3),
    (4, 'This is a great story for kids! It teaches them about accepting everyone, no matter what.', 1),
    (5, 'Mads Mikkelsen is truly one of the most talented actors of our time. He brings so much depth to his characters.', 3),
    (1, 'Búl ertegi degi keıı́pkerlerdiń bári qyzyqty eken! Ásırese, ol Joker Qyz, ol jerde ne ı́steıp jú́r? 😂', 2),
    (1, 'Músafańyń titan búmerangy qandaı qyzyq! Ol onymen qyzyqty ańdardy aulaıtyn shyǵar? 🤔', 3),
    (3, 'Kóktemniń qudaıy men jer qozysyńyń mahabbat tarihy? Ne degen ádemi! 🤩', 1),
    (4, 'Árine, Sparky basqalardan eńreselenǵisi kelse, oǵan eshkimniń eshkanndaı qatysy joq! 🌈', 3),
    (5, 'Ol "The Kings Speech" filmindegi Tristan rólin somdaǵan ba? 😉', 3),
    (5, 'I agree! He can play anything and be absolutely believable. Hes a chameleon of the acting world!', 1);

INSERT INTO likes (user_id, post_id) VALUES
    (1, 1),
    (2, 2),
    (3, 3),
    (1, 3),
    (2, 3),
    (1, 4),
    (2, 5);

INSERT INTO dislikes (user_id, post_id) VALUES
    (2, 1),
    (3, 2);
